package seed

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"yoink/crawler"
	"yoink/utils/env"
	"yoink/utils/logging"
	mysqs "yoink/utils/myaws/sqs"
	myredis "yoink/utils/redis"
	"yoink/utils/resend"
)

// set to true when the aim is to have a million messages in sqs, and false after that stage
var IsDiscovering bool = true

func Crawler() {
	// fetch crawler queue url
	sqsUrl, err := mysqs.GetQueueURL(env.EnvValue.CrawlerSqsName)
	if err != nil {
		fmt.Printf("error in fetching queue url --- %s", err.Error())
		return
	}

	// set up periodic logging
	if env.ConfigValue.Application == "dev" {
		logging.StartHeartbeat(&myredis.CacheHit, &myredis.CacheMiss)
		logging.SendHearbeatMailCrawler(&myredis.CacheHit, &myredis.CacheMiss)
		fmt.Println("Logging and mail services started")
	}

	// initialize config values
	Workers, err := strconv.Atoi(env.ConfigValue.Workers)
	if err != nil {
		fmt.Println("cant convert workers to int")
		return
	}
	NoOfSQSMessages, err := strconv.Atoi(env.ConfigValue.NoOfSQSMessages)
	if err != nil {
		fmt.Println("cant convert no of sqs messages to int")
		return
	}
	var wg sync.WaitGroup

	// queue monitor setup
	if err := mysqs.GetNoOfMessages(sqsUrl); err != nil {
		panic(fmt.Errorf("error in getting messages in queue --- %s", err.Error()))
	}
	mysqs.StartQueueMonitor(sqsUrl)

	// main logic
	t1 := time.Now().UnixMilli()
	for i := 0; i < Workers; i++ {
		task := func() {
			for {
				if IsDiscovering && atomic.LoadInt64(&mysqs.NoOfSQSMessages) >= int64(NoOfSQSMessages) {
					return
				}
				if atomic.LoadInt64(&mysqs.NoOfSQSMessages) <= 0 {
					fmt.Println("empty sqs")
					return
				}
				t1 := time.Now().UnixMilli()
				if err := crawler.Crawl(IsDiscovering, sqsUrl); err != nil {
					fmt.Println("error in main.go => ", err)
					continue
				}
				t2 := time.Now().UnixMilli()
				fmt.Println("time = ", t2-t1, " ms")
			}
		}
		fmt.Printf("worker %d started\n", i+1)
		wg.Go(task)
	}

	// wait for all 10 to finish
	wg.Wait()

	t2 := time.Now().UnixMilli()

	// summary mail
	resend.SendEmail(
		fmt.Sprintf(`
		====== SUMMARY ======
		urls discovered: %d,
		cache hits: %d,
		cache misses:%d,
		workers:%d,
		runtime:%d,
		`, atomic.LoadInt64(&mysqs.NoOfSQSMessages),
			atomic.LoadInt64(&myredis.CacheHit),
			atomic.LoadInt64(&myredis.CacheMiss),
			Workers,
			t2-t1,
		),
		"COMPLETED EC2 CRAWLING",
	)
}
