package config

import (
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
)

type CErrorOptions struct {
	User  *sentry.User
	Err   string
	Level sentry.Level
}

func CustomError(options CErrorOptions) {
	hub := sentry.CurrentHub()
	if options.User != nil {
		hub.Scope().SetUser(*options.User)
	}
	hub.Scope().SetLevel(options.Level)
	hub.CaptureMessage(options.Err)
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		fmt.Printf("%+v", options)
	}
}

type ResponseError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    *any   `json:"data"`
}
type ResponseSuccess struct {
	Status  bool
	Message string
}
