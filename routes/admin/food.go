package admin

import (
	"time"

	"github.com/abasalt-yar/college-food-reserver/body"
	"github.com/abasalt-yar/college-food-reserver/config"
	"github.com/abasalt-yar/college-food-reserver/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddFood(c *gin.Context) {
	db := c.Keys["PSQLDB"].(*gorm.DB)
	admin := c.Keys["Admin"].(*models.Admin)
	if admin.Position != "COLLEGE_ADMIN" && admin.Position != "COLLEGE_OWNER" {
		c.JSON(403, &config.ResponseError{Status: false, Message: "AUTHORIZATION_FAILED"})
		return
	}
	var body body.AdminAddFoodBody
	if err := body.Validate(c); err != nil {
		c.JSON(400, err)
		return
	}
	inst := db.Table("foods").Create(&models.Food{FoodName: body.FoodName, Expensive: body.IsExpensive, CollegeId: admin.CollegeId, AddedBy: admin.ID, CreatedDate: time.Now()})
	if err := inst.Error; err != nil {
		c.JSON(500, &config.ResponseError{Status: false, Message: "INTERNAL_SERVER_ERROR"})
		return
	}
	c.JSON(201, &config.ResponseSuccess{Status: true, Message: "FOOD_ADDED"})

}

func RemoveFood(c *gin.Context) {
	db := c.Keys["PSQLDB"].(*gorm.DB)
	admin := c.Keys["Admin"].(*models.Admin)
	if admin.Position != "COLLEGE_ADMIN" && admin.Position != "COLLEGE_OWNER" {
		c.JSON(403, &config.ResponseError{Status: false, Message: "AUTHORIZATION_FAILED"})
		return
	}
	var body body.AdminRemoveFood
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"})
		return
	}
	del := db.Table("foods").Delete(&models.Food{ID: body.FoodID}).Limit(1)
	if del.Error != nil || del.RowsAffected == 0 {
		c.JSON(404, &config.ResponseError{Status: false, Message: "FOOD_NOT_FOUND"})
		return
	}
	c.JSON(200, &config.ResponseSuccess{Status: true, Message: "FOOD_DELETED"})
}
