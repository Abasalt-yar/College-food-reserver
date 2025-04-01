package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func connectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return rdb
}
