package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"yoink/utils/env"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func NewClient() error{
	redis_db, err := strconv.Atoi(env.EnvValue.RedisDatabase)
	if err != nil{
		return err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     env.EnvValue.RedisAddress,
		Username: env.EnvValue.RedisUsername,
		Password: env.EnvValue.RedisPassword,
		DB:       redis_db,
	})

	RDB = rdb
	fmt.Println("Redis database connection success")
	return nil
}

func SetCache(key string, value string) error {
	return RDB.Set(context.Background(), key, value, 7*24*time.Hour).Err()
}

func GetCache(key string) (string, error) {
	return RDB.Get(context.Background(), key).Result()
}