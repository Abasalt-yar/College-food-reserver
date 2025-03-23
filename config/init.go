package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ConfigInitilization struct {
	PSQLDB *gorm.DB
	REDIS  *redis.Client
}

func Init() *ConfigInitilization {
	return &ConfigInitilization{
		PSQLDB: connectPSQLDatabase(),
		REDIS:  connectRedis(),
	}
}
