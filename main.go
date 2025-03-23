package main

import (
	"github.com/abasalt-yar/college-food-reserver/config"
	"github.com/abasalt-yar/college-food-reserver/routes/admin"
	"github.com/abasalt-yar/college-food-reserver/routes/middleware"
	"github.com/abasalt-yar/college-food-reserver/routes/student"
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config_data := config.Init()
	main_router := gin.New()
	var traceSampleRate float64 = .5
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		main_router.Use(gin.Logger())
		println("LOG LEVEL DEBUG")
		traceSampleRate = 1.0
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://66ed465f1236561da81dbc733e31aa60@o4508252312109056.ingest.de.sentry.io/4508252315123792",
		EnableTracing:    true,
		TracesSampleRate: traceSampleRate,
	}); err != nil {
		fmt.Printf("Sentry Inilitization Failed: %v\n", err)
	}
	main_router.Use(middleware.CORSMiddleware())

	main_router.Use(sentrygin.New(sentrygin.Options{}))

	main_router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "Pong!")
	})

	main_router.Use(func(ctx *gin.Context) {
		ctx.Set("PSQLDB", config_data.PSQLDB)
		ctx.Set("REDIS", config_data.REDIS)
	})
	student.Init(main_router.Group("/student"))
	admin.Init(main_router.Group("/admin"))

	main_router.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
