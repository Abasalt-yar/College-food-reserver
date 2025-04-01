package admin

import (
	"errors"
	"time"

	"github.com/abasalt-yar/college-food-reserver/body"
	"github.com/abasalt-yar/college-food-reserver/config"
	"github.com/abasalt-yar/college-food-reserver/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AddStudent(c *gin.Context) {
	db := c.Keys["PSQLDB"].(*gorm.DB)
	admin := c.Keys["Admin"].(*models.Admin)
	if admin.Position != "COLLEGE_ADMIN" && admin.Position != "COLLEGE_OWNER" {
		c.JSON(403, &config.ResponseError{Status: false, Message: "AUTHORIZATION_FAILED"})
		return
	}
	var body body.AdminAddStudent
	if err := body.Validate(c); err != nil {
		c.JSON(400, err)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(422, &config.ResponseError{Status: false, Message: "INVALID_PASSWORD"})
		return
	}
	inst := db.Table("students").Create(&models.Student{Username: body.Username, Password: string(hash), FirstName: body.FirstName, LastName: body.LastName, PhoneNumber: body.PhoneNumber, Position: body.Position, CollegeId: admin.CollegeId, AddedBy: admin.ID, CreatedDate: time.Now()})
	if err := inst.Error; err != nil {

		if errors.Is(err, gorm.ErrCheckConstraintViolated) || errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrInvalidValueOfLength) {
			c.JSON(400, &config.ResponseError{Status: false, Message: "DATA_CHECK_FAILED"})
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(400, &config.ResponseError{Status: false, Message: "DUPLICATE"})
		} else {
			c.JSON(500, &config.ResponseError{Status: false, Message: "INTERNAL_SERVER_ERROR"})
		}
		return
	}
	c.JSON(201, &config.ResponseSuccess{Status: true, Message: "STUDENT_ADDED"})
}
func RemoveStudent(c *gin.Context) {
	db := c.Keys["PSQLDB"].(*gorm.DB)
	admin := c.Keys["Admin"].(*models.Admin)
	if admin.Position != "COLLEGE_ADMIN" && admin.Position != "COLLEGE_OWNER" {
		c.JSON(403, &config.ResponseError{Status: false, Message: "AUTHORIZATION_FAILED"})
		return
	}
	var body body.AdminRemoveStudent
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"})
		return
	}
	del := db.Table("students").Delete(&models.Student{ID: body.StudentID, CollegeId: admin.CollegeId}).Limit(1)
	if del.Error != nil {
		c.JSON(404, &config.ResponseError{Status: false, Message: "STUDENT_NOT_FOUND"})
		return
	}
	c.JSON(200, &config.ResponseSuccess{Status: true, Message: "STUDENT_DELETED"})
}

func UpdateStudentPosition(c *gin.Context) {
	db := c.Keys["PSQLDB"].(*gorm.DB)
	admin := c.Keys["Admin"].(*models.Admin)
	if admin.Position != "COLLEGE_ADMIN" && admin.Position != "COLLEGE_OWNER" {
		c.JSON(403, &config.ResponseError{Status: false, Message: "AUTHORIZATION_FAILED"})
		return
	}
	var body body.AdminUpdateStudentPosition
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"})
		return
	}
	del := db.Table("students").Where(&models.Student{ID: body.StudentID, CollegeId: admin.CollegeId}).UpdateColumn("position", body.Position).Limit(1)
	if del.Error != nil || del.RowsAffected == 0 {
		c.JSON(404, &config.ResponseError{Status: false, Message: "STUDENT_NOT_FOUND"})
		return
	}
	c.JSON(200, &config.ResponseSuccess{Status: true, Message: "STUDENT_UPDATED"})
}
