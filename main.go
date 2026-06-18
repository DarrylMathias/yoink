package main

import (
	"fmt"
	"time"
	"yoink/app"
	// "yoink/crawler"
)

func main() {
	t1 := time.Now().UnixMilli()
	app.App()
	// if err := crawler.Crawl(); err != nil{
	// 	panic(err)
	// }
	t2 := time.Now().UnixMilli()
	fmt.Println("round trip time (ms)", t2-t1)
}
