package app

import (
	"fmt"

	"yoink/utils/database"
	"yoink/utils/env"
	"yoink/utils/myaws"
	"yoink/utils/myaws/s3"
	mysqs "yoink/utils/myaws/sqs"
	"yoink/utils/redis"
	"yoink/utils/resend"
	"yoink/utils/upstash"
)

func App(){
	err := env.NewEnv()
	if err != nil {
		panic(fmt.Errorf("error in parsing env --- %s", err.Error()))
	}
	err = myaws.GetConfig()
	if err != nil {
		panic(fmt.Errorf("error in aws config --- %s", err.Error()))
	}

	mysqs.GetSQSClient()
	s3.GetS3Client()
	resend.GetResendClient()

	database.NewDatabase()
	err = redis.NewClient()
	if err != nil{
		panic(fmt.Errorf("error in redis config --- %s", err.Error()))
	}
	err = upstash.NewClient()
	if err != nil{
		panic(fmt.Errorf("error in upstash config --- %s", err.Error()))
	}
}