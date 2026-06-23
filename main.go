package main

import (
	"yoink/app"
	// "yoink/indexer"
	// "yoink/utils/env"
	"yoink/seed"
)

func main() {
    app.App()
	// seed.SeedSQS()
	// seed.Crawler()
	// if err := indexer.Indexer(&env.EnvValue.IndexerSqsName); err != nil{
	// 	panic(err)
	// }
	seed.IndexerSeed()
}
