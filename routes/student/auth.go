package student

import (
	"os"
	"strconv"
	"time"

	"github.com/abasalt-yar/college-food-reserver/body"
	"github.com/abasalt-yar/college-food-reserver/config"
	"github.com/abasalt-yar/college-food-reserver/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func StudentLogin(c *gin.Context) {
	db := (c.Keys["PSQLDB"].(*gorm.DB))
	student := models.Student{}
	var requestBody body.StudentLoginBody
	if err := requestBody.Validate(c); err != nil {
		c.JSON(400, err)
		return
	}

	if err := db.First(&student, &models.Student{Username: requestBody.Username}).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "یوزرنیم و یا پسورد اشتباه میباشد",
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
	if bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(requestBody.Password)) != nil {
		c.JSON(401, config.ResponseError{
			Status:  false,
			Message: "یوزرنیم و یا پسورد اشتباه میباشد",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":         os.Getenv("JWT_ISS"),
		"nbf":         time.Now().Unix(),
		"exp":         time.Now().AddDate(0, 0, 30).Unix(),
		"sub":         strconv.FormatUint(uint64(student.ID), 10),
		"username":    student.Username,
		"incremental": student.Counter,
		"admin":       false,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, config.ResponseError{
			Status:  false,
			Message: "SERVER_INTERNAL_ERROR",
		})
		return
	}
	c.JSON(200, body.LoginTokenResponse{
		Status:      true,
		AccessToken: tokenString,
	})

}

func StudentChangePassword(c *gin.Context) {
	var requestBody body.StudentCPassBody
	if er := requestBody.Validate(c); er != nil {
		c.JSON(422, er)
		return
	}
	db := c.Keys["PSQLDB"].(*gorm.DB)
	user := c.Keys["Student"].(*models.Student)
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Current)) != nil {
		c.JSON(401, &config.ResponseError{Status: false, Message: "CURRENT_PASSWORD_INCORRECT"})
		return
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(requestBody.New), 10)
	if err != nil {
		c.JSON(422, &config.ResponseError{Status: false, Message: "INVALID_PASSWORD"})
		return
	}
	if er := db.Model(&models.Student{}).Limit(1).Where("id = ?", user.ID).UpdateColumns(&models.Student{Password: string(newHash), Counter: user.Counter + 1}).Error; er != nil {
		c.JSON(500, &config.ResponseError{Status: false, Message: "INTERNAL_SERVER_ERROR"})
		return
	}
	c.JSON(200, &body.StudentCPassResponse{Status: true, Message: "پسورد شما با موفقیت تغییر کرد", Logout: true})
}
