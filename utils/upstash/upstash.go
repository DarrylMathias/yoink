package upstash

import (
	"context"
	"fmt"
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

  return nil
}