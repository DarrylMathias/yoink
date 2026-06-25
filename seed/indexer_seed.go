package seed

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"yoink/indexer/store"
	"yoink/indexer"
	"yoink/utils/env"
	"yoink/utils/logging"
	mysqs "yoink/utils/myaws/sqs"
	"yoink/utils/resend"
)

func IndexerSeed(){
	// fetch crawler queue url
	sqsUrl, err := mysqs.GetQueueURL(env.EnvValue.IndexerSqsName)
	if err != nil{
		fmt.Printf("error in fetching queue url --- %s", err.Error())
		return
	}
	
	var indexedPages int64

	// set up periodic logging
	logging.StartHeartbeatIndexer(&indexedPages)
	logging.SendHearbeatMailIndexer(&indexedPages)
	fmt.Println("Logging and mail services started")

	workers, err := strconv.Atoi(env.ConfigValue.Workers)
	if err != nil{
		panic(err)
	}
	var wg sync.WaitGroup

	// queue monitor setup
	if err := mysqs.GetNoOfMessages(sqsUrl); err != nil{
		panic(fmt.Errorf("error in getting messages in queue --- %s", err.Error()))
	}
	mysqs.StartQueueMonitor(sqsUrl)

	store.Init()
	t1 := time.Now().UnixMilli()
	for w:=0; w<workers; w++{
		fmt.Println("started worker", w+1)
		wg.Go(func() {
			for{
				count, err := indexer.Indexer(sqsUrl)
				if errors.Is(err, indexer.ErrEmptyQueue){
					fmt.Println("worker exiting: empty queue")
					return
				}
				if err != nil{
					fmt.Println("indexer error:", err)
					continue
				}
				if count > 0 {
					atomic.AddInt64(&indexedPages, int64(count))
				}
			}
		})
	}
	wg.Wait()
	t2 := time.Now().UnixMilli()

	// final flush to ensure remaining data is saved
	fmt.Println("flushing remaining data to disk...")
	if err := store.Flush(); err != nil {
		fmt.Printf("final flush error: %v\n", err)
	} else {
		fmt.Println("final flush successful")
	}
	
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