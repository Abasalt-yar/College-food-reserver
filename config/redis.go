package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func connectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		// Addr: os.Getenv("REDIS_ADDR"),
		Addr: "redis-15403.c294.ap-northeast-1-2.ec2.redns.redis-cloud.com:15403",
		Username: "default",
		Password: "GiSkaAOYY26Zoah3E8MzH0f9UKxw6Kat",
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return rdb
}
