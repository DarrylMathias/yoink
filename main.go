package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"yoink/app"
	"yoink/indexer"
	"yoink/utils/env"
	"yoink/utils/logging"
	mysqs "yoink/utils/myaws/sqs"
	"yoink/utils/resend"
	// "yoink/seed"
)

func main() {
    app.App()
	// seed.SeedSQS()
	// seed.Crawler()
	// seed.IndexerSeed()

	// fetch crawler queue url
	sqsUrl, err := mysqs.GetQueueURL(env.EnvValue.IndexerSqsName)
	if err != nil{
		fmt.Printf("error in fetching queue url --- %s", err.Error())
		return
	}
	
	// set up periodic logging
	if env.ConfigValue.Application == "dev"{
		logging.StartHeartbeatIndexer()
		logging.SendHearbeatMailIndexer()
		fmt.Println("Logging and mail services started")
	}

	workers, err := strconv.Atoi(env.ConfigValue.Workers)
	if err != nil{
		panic(err)
	}
	var wg sync.WaitGroup

	t1 := time.Now().UnixMilli()
	for w:=0; w<workers; w++{
		fmt.Println("started worker", w+1)
		wg.Go(func() {
			for{
				err := indexer.Indexer(sqsUrl)
				if errors.Is(err, indexer.ErrEmptyQueue){
					return
				}
				if err != nil{
					fmt.Println("indexer error:", err)
					continue
				}
			}
		})
	}
	t2 := time.Now().UnixMilli()
	wg.Wait()
	
	// summary mail
	resend.SendEmail(
		fmt.Sprintf(`
		====== SUMMARY ======
		urls discovered: %d,
		workers:%d,
		runtime:%d mins,
		`, 	atomic.LoadInt64(&mysqs.NoOfSQSMessages),
			workers,
			(t2-t1)/(1000*60),
		),
		"COMPLETED EC2 INDEXING",
	)
}
