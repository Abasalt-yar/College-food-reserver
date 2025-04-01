package student

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/abasalt-yar/college-food-reserver/body"
	"github.com/abasalt-yar/college-food-reserver/config"
	"github.com/abasalt-yar/college-food-reserver/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func StudentGetMe(c *gin.Context) {
	var hi struct {
		Status  bool
		Message string
	}
	user := c.Keys["Student"].(*models.Student)
	hi.Status = true
	var mySecondPass string
	if user.SecondPass == nil {
		mySecondPass = "NO_SECOND_PASS"
	} else {
		mySecondPass = *user.SecondPass
	}
	hi.Message = fmt.Sprintf("Hi %s | Your Second Pass = %s", user.FirstName, mySecondPass)
	c.JSON(200, hi)
}

func StudentTransferBalance(c *gin.Context) {
	user := c.Keys["Student"].(*models.Student)
	var requestBody body.StudentTBalanceBody
	if err := requestBody.Validate(c); err != nil {
		c.JSON(422, err)
		return
	}
	if user.Username == requestBody.Target {
		c.JSON(401, config.ResponseError{Status: false, Message: "CANT_DO"})
		return
	}
	if user.Balance < int16(requestBody.Balance) {
		c.JSON(401, config.ResponseError{Status: false, Message: "INSUFFICIENT_BALANCE"})
		return
	}
	db := c.Keys["PSQLDB"].(*gorm.DB)
	var target models.Student
	if err := db.Select("balance", "id", "first_name", "last_name").First(&target, &models.Student{Username: requestBody.Target}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "TARGET_NOT_FOUND",
			})
			return
		}
		go config.CustomError(config.CErrorOptions{
			Err:   err.Error(),
			Level: sentry.LevelFatal,
		})
		c.JSON(500, config.ResponseError{
			Status:  false,
			Message: "SERVER_INTERNAL_ERROR",
		})
		return
	}
	if err := db.Model(&models.Student{}).Limit(1).Where("id = ?", user.ID).UpdateColumn("balance", gorm.Expr("balance - ?", requestBody.Balance)).Error; err != nil {
		c.JSON(422, config.ResponseError{Status: false, Message: "INSUFFICIENT_BALANCE"})
		return
	}
	if err := db.Model(&models.Student{}).Limit(1).Where("id = ?", target.ID).UpdateColumn("balance", gorm.Expr("balance + ?", requestBody.Balance)).Error; err != nil {
		go config.CustomError(config.CErrorOptions{
			Err:   err.Error(),
			Level: sentry.LevelFatal,
			User:  &sentry.User{Username: user.Username, ID: fmt.Sprintf("%d", target.ID)},
		})
		c.JSON(500, config.ResponseError{
			Status:  false,
			Message: "SERVER_INTERNAL_ERROR",
		})
		return
	}
	db.Table("btransfer_log").Create(&models.Btransferlog{
		FromWho:     user.ID,
		ToWho:       target.ID,
		Price:       uint(requestBody.Balance),
		CreatedDate: time.Now(),
	})
	c.JSON(200, body.StudentTBalanceResponse{
		Status:     true,
		To:         fmt.Sprintf("%s %s", target.FirstName, target.LastName),
		NewBalance: user.Balance - int16(requestBody.Balance),
	})

}

func StudentChange2Pass(c *gin.Context) {
	var requestBody body.StudentC2Pass
	if err := requestBody.Validate(c); err != nil {
		c.JSON(422, err)
		return
	}
	user := c.Keys["Student"].(*models.Student)
	db := c.Keys["PSQLDB"].(*gorm.DB)
	if user.SecondPass != nil && user.SecondPass == &requestBody.New {
		c.JSON(200, &body.StudentCPassResponse{Status: true, Message: "پسورد شما با موفقیت تغییر یافت", Logout: false})
		return
	}
	if err := db.Model(&models.Student{}).Limit(1).Where("id = ?", user.ID).UpdateColumn("second_pass", requestBody.New).Error; err != nil {
		c.JSON(500, &config.ResponseError{Status: false, Message: "INTERNAL_SERVER_ERROR"})
		return
	}
	c.JSON(200, &body.StudentCPassResponse{Status: true, Message: "پسورد شما با موفقیت تغییر یافت", Logout: false})

}

func RequestAddBalance(c *gin.Context) {
	howMuch, err := strconv.ParseUint(c.Query("price"), 10, 64)
	if err != nil || howMuch%1000 != 0 || howMuch < 1000 || howMuch > 200_000 {
		c.JSON(422, &config.ResponseError{Status: false, Message: "PRICE_PARAMETER_NOT_VALID"})
		return
	}

	user := c.Keys["Student"].(*models.Student)
	//* Fix: Get Payment URL
	var resp config.ZarinpalRequest
	if resp.DoRequest(fmt.Sprintf(`{"merchant_id": "%s","amount": %d,"description": "%s","callback_url": "%s","metadata": [{"mobile": "%s"}]}`, os.Getenv("ZARINPAL_MERCHANT_ID"), howMuch*10, fmt.Sprintf("ADDBALANCE-%d", user.ID), "https://google.com/", user.PhoneNumber)) != nil {
		c.JSON(500, &config.ResponseError{Status: false, Message: "INTERNAL_SERVER_ERROR"})
		return
	}
	if resp.Data.Code != 100 {
		c.JSON(500, &config.ResponseError{Status: false, Message: "NOT_VALID_RESPONSE_FROM_GATEWAY"})
		return
	}
	rdb := c.Keys["REDIS"].(*redis.Client)
	if rdb.JSONSet(context.Background(), fmt.Sprintf(`CFR:AB#%s`, resp.Data.Authority), "$", config.ZarinaplRedis{
		UserID: user.ID,
		Amount: howMuch,
	}).Err() != nil {
		c.JSON(500, &config.ResponseError{Status: false, Message: "FAILED_TO_SAVE_PAYMENT"})
		return
	}
	rdb.Expire(context.Background(), fmt.Sprintf(`CFR:AB#%s`, resp.Data.Authority), 2*time.Hour)
	c.JSON(201, &body.StudentAddBalanceResponse{Status: true, Authority: resp.Data.Authority})
}

func VerifyPayment(c *gin.Context) {
	var requestBody body.StudentVerifyPayment
	if err := requestBody.Validate(c); err != nil {
		c.JSON(422, err)
		return
	}
	var verified config.ZarinpalVerify
	rdb := c.Keys["REDIS"].(*redis.Client)
	rdbKey := fmt.Sprintf("CFR:AB#%s", requestBody.Authority)
	result, er := rdb.JSONGet(context.Background(), rdbKey, "$").Result()
	var zredis []config.ZarinaplRedis
	if er != nil || len(result) < 1 {
		c.JSON(400, &config.ResponseError{Status: false, Message: ".PAYMENT_NOT_FOUND"})
		return
	}
	if err := json.Unmarshal([]byte(result), &zredis); err != nil || len(zredis) != 1 {
		c.JSON(400, &config.ResponseError{Status: false, Message: "PAYMENT_NOT_FOUND."})
		return
	}

	if err := verified.DoRequest(fmt.Sprintf(`{"merchant_id": "%s","amount": %d,"authority": "%s"}`, os.Getenv("ZARINPAL_MERCHANT_ID"), zredis[0].Amount*10, requestBody.Authority)); err != nil {
		if err.Error() == "PAYMENT_NOT_COMPLETE" {
			c.JSON(400, &config.ResponseError{Status: false, Message: "PAYMENT_NOT_COMPLETE"})
		} else {
			c.JSON(500, &config.ResponseError{Status: false, Message: "INTERNAL_SERVER_ERROR"})
		}
		return
	}
	if verified.Data.Code == 101 {
		c.JSON(400, &config.ResponseError{Status: false, Message: "PAYMENT_ALREADY_VERIFIED"})
		rdb.JSONDel(context.Background(), rdbKey, "$")
		return
	}
	if verified.Data.Code != 100 {
		c.JSON(400, &config.ResponseError{Status: false, Message: "PAYMENT_NOT_COMPLETE"})
		return
	}
	rdb.JSONDel(context.Background(), rdbKey, "$")
	db := c.Keys["PSQLDB"].(*gorm.DB)

	if insert := db.Table("btransactions").Create(&models.Btransaction{
		Authority:   requestBody.Authority,
		RefID:       uint(verified.Data.RefID),
		PayedBy:     zredis[0].UserID,
		Price:       uint(zredis[0].Amount),
		CreatedDate: time.Now(),
	}); insert.Error != nil {
		go config.CustomError(config.CErrorOptions{
			Err:   insert.Error.Error(),
			Level: sentry.LevelFatal,
			User:  &sentry.User{ID: fmt.Sprintf("%d", zredis[0].UserID), Data: map[string]string{"authority": requestBody.Authority}},
		})
		c.JSON(500, config.ResponseError{
			Status:  false,
			Message: "SERVER_INTERNAL_ERROR",
		})
		return
	}
	if insert := db.Model(&models.Student{}).Limit(1).Where("id = ?", zredis[0].UserID).UpdateColumn("balance", gorm.Expr("balance + ?", zredis[0].Amount)); insert.Error != nil {
		go config.CustomError(config.CErrorOptions{
			Err:   insert.Error.Error(),
			Level: sentry.LevelFatal,
			User:  &sentry.User{ID: fmt.Sprintf("%d", zredis[0].UserID), Data: map[string]string{"authority": requestBody.Authority}},
		})
		c.JSON(500, config.ResponseError{
			Status:  false,
			Message: "SERVER_INTERNAL_ERROR",
		})
		return
	}
	c.JSON(200, &body.StudentVerifyPaymentResponse{Status: true, RefID: verified.Data.RefID, Price: zredis[0].Amount})
}
