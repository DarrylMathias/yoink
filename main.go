package main

import (
	"fmt"
	"yoink/app"

	"yoink/ranking"
)


func main() {
    app.App()
	// seed.SeedSQS()
	// seed.Crawler()
	//seed.IndexerSeedSQS()
	// seed.IndexerSeed()
	ranking.Init()
	results, err := ranking.RankPages("why are apples red?")
	if err != nil{
		panic(err)
	}
	for i, result := range results{
		fmt.Printf("%d. %s => %s\n", i, result.Title, result.Url)
	}
}
