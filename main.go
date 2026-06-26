package main

import (
	"yoink/app"
	"yoink/ranking"
	"fmt"
	// "yoink/ranking/fetch"
)

func main() {
    app.App()
	// seed.SeedSQS()
	// seed.Crawler()
	//seed.IndexerSeedSQS()
	// seed.IndexerSeed()
	results, err := ranking.RankPages("why are apples red?")
	if err != nil{
		panic(err)
	}
	fmt.Println(results)
}
