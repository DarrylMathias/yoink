package main

import (
	"yoink/app"
	// "yoink/crawler/extract/hashtable"
	"yoink/crawler/extract"
)

func main() {
	app.App()
	err := extract.ExtractURLData()
	if err != nil {
		panic(err)
	}
	// _, err := hashtable.AlreadySeen("abd")
	// if err != nil {
	// 	panic(err)
	// }
}
