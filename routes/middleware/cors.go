package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cors_origin string = "*"

		if os.Getenv("PROJECT_LEVEL") == "test" {
			cors_origin = "*"
		} else {
			if c.GetHeader("origin") == "http://www.samad.app" {
				cors_origin = "https://www.samad.app"
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", cors_origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
