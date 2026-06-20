package main

import (
	// "fmt"
	// "sync"
	// "sync/atomic"
	// "time"
	"yoink/app"
	"yoink/seed"
	// "yoink/crawler"
	// mysqs "yoink/utils/myaws/sqs"
	// "yoink/utils/resend"
)

func main() {
    app.App()
	// const Workers = 50
	// var wg sync.WaitGroup

	// t1 := time.Now().UnixMilli()
	// for i := 0; i < Workers; i++ {
	// 	task := func() {
	// 		for {
	// 			if app.IsDiscovering && atomic.LoadInt64(&mysqs.NoOfSQSMessages) >= 1_0_000 {
	// 				return
	// 			}
	// 			if !app.IsDiscovering && atomic.LoadInt64(&mysqs.NoOfSQSMessages )<= 0{
	// 				return
	// 			}
	// 			t1 := time.Now().UnixMilli()
	// 			if err := crawler.Crawl(); err != nil{
	// 				fmt.Println("error in main.go => ", err)
	// 				continue
	// 			}
	// 			t2 := time.Now().UnixMilli()
	// 			fmt.Println("time = ", t2-t1, " ms")
	// 		}
	// 	}
	// 	fmt.Printf("worker %d started\n", i+1)
	// 	wg.Go(task)
	// }
	// // wait for all 10 to finish
	// wg.Wait()

	// t2 := time.Now().UnixMilli()

	// resend.SendEmail(
	// 	fmt.Sprintf(`
	// 	====== SUMMARY ======
	// 	urls discovered: %d,
	// 	cache hits: %d,
	// 	cache misses:%d,
	// 	workers:%d,
	// 	runtime:%d,
	// 	`, 	atomic.LoadInt64(&mysqs.NoOfSQSMessages),
	// 		atomic.LoadInt64(&app.CacheHit),
	// 		atomic.LoadInt64(&app.CacheMiss),
	// 		Workers,
	// 		t2-t1,
	// 	),
	// 	"COMPLETED EC2 CRAWLING",
	// )
	seed.SeedSQS()
}
