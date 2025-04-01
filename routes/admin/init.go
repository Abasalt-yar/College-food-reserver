package admin

import (
	"github.com/abasalt-yar/college-food-reserver/routes/middleware"

	"github.com/gin-gonic/gin"
)

func Init(rg *gin.RouterGroup) {
	rg.POST("/login", AdminLogin)
	rg.PUT("/changepass", middleware.AuthMiddlewareAdmin, AdminChangePassword)
	rg.POST("/food/add", middleware.AuthMiddlewareAdmin, AddFood)
	rg.DELETE("/food/remove", middleware.AuthMiddlewareAdmin, RemoveFood)
	rg.POST("/students/add", middleware.AuthMiddlewareAdmin, AddStudent)
	rg.DELETE("/students/remove", middleware.AuthMiddlewareAdmin, RemoveStudent)
	rg.PUT("/students/updateposition", middleware.AuthMiddlewareAdmin, UpdateStudentPosition)
}
