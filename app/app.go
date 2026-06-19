package app

import (
	"fmt"
	"sync/atomic"
	"time"

	"yoink/utils/database"
	"yoink/utils/env"
	"yoink/utils/myaws"
	"yoink/utils/myaws/s3"
	mysqs "yoink/utils/myaws/sqs"
	"yoink/utils/redis"
	"yoink/utils/resend"
	"yoink/utils/upstash"
)

func StartHeartbeat() {
	go func() {
		for {
			fmt.Printf(
				"[HEARTBEAT] frontier=%d hits=%d misses=%d\n",
				atomic.LoadInt64(&mysqs.NoOfSQSMessages),
				atomic.LoadInt64(&CacheHit),
				atomic.LoadInt64(&CacheMiss),
			)

			time.Sleep(time.Minute)
		}
	}()
}

func SendHearbeatMail() {
	go func() {
		for {
			err, mailId := resend.SendEmail(
				fmt.Sprintf(
					"[HEARTBEAT] frontier=%d hits=%d misses=%d\n",
					atomic.LoadInt64(&mysqs.NoOfSQSMessages),
					atomic.LoadInt64(&CacheHit),
					atomic.LoadInt64(&CacheMiss),
				),
				"Crawling updates",
			)

			if err != nil {
				fmt.Printf("heartbeat mail failed: %v\n", err)
			} else {
				fmt.Printf("heartbeat mail sent: %s\n", mailId)
			}

			time.Sleep(6 * time.Hour)
		}
	}()
}

var CacheHit int64 = 0
var CacheMiss int64 = 0

// set to true when the aim is to have a million messages in sqs, and false after that stage
var IsDiscovering bool = false

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
	resend.GetResendClient()

	database.NewDatabase(env.EnvValue)
	err = redis.NewClient()
	if err != nil{
		panic(fmt.Errorf("error in redis config --- %s", err.Error()))
	}
	err = upstash.NewClient()
	if err != nil{
		panic(fmt.Errorf("error in upstash config --- %s", err.Error()))
	}

	// periodic logging
	StartHeartbeat()
	SendHearbeatMail()
	fmt.Println("Logging and mail services started")
	
}