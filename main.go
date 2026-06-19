package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"yoink/app"
	"yoink/crawler"
)

func main() {
    app.App()
	const Workers = 10
	var wg sync.WaitGroup

	t1 := time.Now().UnixMilli()
	for i := 0; i < Workers; i++ {
		task := func() {
			for atomic.LoadInt64(&app.Counter) < 100 {
				fmt.Println("=============================================================")
				t1 := time.Now().UnixMilli()
				if err := crawler.Crawl(); err != nil{
					fmt.Println("error in main.go => ", err)
					continue
				}
				t2 := time.Now().UnixMilli()
				fmt.Println("time = ", t2-t1, " ms")
				fmt.Println("=============================================================")
			}
		}
		fmt.Printf("worker %d started\n", i+1)
		wg.Go(task)
	}
	t2 := time.Now().UnixMilli()
	// wait for all 10 to finish
	wg.Wait()
	fmt.Println("total time = ", t2-t1, " ms")
	fmt.Println("hit/miss ration = ", float64(atomic.LoadInt64(&app.CacheMiss))/float64(atomic.LoadInt64(&app.CacheMiss)))
}
