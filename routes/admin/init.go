package admin

import (
	"github.com/abasalt-yar/college-food-reserver/routes/middleware"

	"github.com/gin-gonic/gin"
)

func Init(rg *gin.RouterGroup) {
	rg.POST("/login", AdminLogin)
	rg.PUT("/changepass", middleware.AuthMiddlewareAdmin, AdminChangePassword)
	rg.POST("/food/add", middleware.AuthMiddlewareAdmin, AddFood)
}
