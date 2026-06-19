package app

import (
	"fmt"

	"yoink/utils/database"
	"yoink/utils/env"
	"yoink/utils/myaws"
	"yoink/utils/myaws/s3"
	mysqs "yoink/utils/myaws/sqs"
	"yoink/utils/redis"
	"yoink/utils/upstash"
)

var CacheHit int64 = 0
var CacheMiss int64 = 0

func App(){
	err := env.NewEnv(".env.prod")
	if err != nil {
		panic(fmt.Errorf("error in parsing env --- %s", err.Error()))
	}
	err = myaws.GetConfig()
	if err != nil {
		panic(fmt.Errorf("error in aws config --- %s", err.Error()))
	}

	mysqs.GetSQSClient()
	s3.GetS3Client()

	database.NewDatabase(env.EnvValue)
	err = redis.NewClient()
	if err != nil{
		panic(fmt.Errorf("error in redis config --- %s", err.Error()))
	}
	err = upstash.NewClient()
	if err != nil{
		panic(fmt.Errorf("error in upstash config --- %s", err.Error()))
	}
	
}