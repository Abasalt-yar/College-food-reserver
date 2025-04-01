package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectPSQLDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_PSQL")), &gorm.Config{TranslateError: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}
	db.Exec("SET time zone 'UTC';")
	return db
}
