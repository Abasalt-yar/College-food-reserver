package student

import (
	"github.com/abasalt-yar/college-food-reserver/routes/middleware"

	"github.com/gin-gonic/gin"
)

func Init(rg *gin.RouterGroup) {
	rg.POST("/login", StudentLogin)
	rg.GET("/profile", middleware.AuthMiddlewareStudent, StudentGetMe)
	rg.POST("/transferBalance", middleware.AuthMiddlewareStudent, StudentTransferBalance)
	rg.PUT("/changepass", middleware.AuthMiddlewareStudent, StudentChangePassword)
	rg.PUT("/change2pass", middleware.AuthMiddlewareStudent, StudentChange2Pass)
	rg.GET("/requestBalance", middleware.AuthMiddlewareStudent, RequestAddBalance)
	rg.POST("/verifyPayment", VerifyPayment)
}
