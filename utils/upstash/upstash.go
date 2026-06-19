package upstash

import (
	"context"
	"fmt"
	"time"
	"yoink/utils/env"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var UpstashClient *redis.Client

func NewClient() error{
  opt, err := redis.ParseURL(fmt.Sprintf("rediss://default:%s@%s:6379", env.EnvValue.UpstashRedisRestToken, env.EnvValue.UpstashRedisRestURL))
  if err != nil{
	return err
  }

  rdb := redis.NewClient(opt)
  UpstashClient = rdb

  fmt.Println("Upstash connection success")

  return nil
}

func SetCache(key string, value string) error {
	return UpstashClient.Set(context.Background(), key, value, 7*24*time.Hour).Err()
}

func GetCache(key string) (string, error) {
	return UpstashClient.Get(context.Background(), key).Result()
}