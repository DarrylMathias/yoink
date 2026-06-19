package main

import (
	"fmt"
	"time"
	"yoink/app"
	"yoink/crawler"
)

func main() {
	app.App()
	for app.Counter<1000{
		fmt.Println("=============================================================")
		t1 := time.Now().UnixMilli()
		if err := crawler.Crawl(); err != nil{
			fmt.Println("error in main.go => %w", err)
			continue
		}
		t2 := time.Now().UnixMilli()
		fmt.Println("round trip time (ms)", t2-t1)
		fmt.Println("=============================================================")
	}
}
