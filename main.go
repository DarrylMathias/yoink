package main

import (
	"yoink/app"
	"yoink/seed"
)

func main() {
    app.App()
	// seed.SeedSQS()
	// seed.Crawler()
	//seed.IndexerSeedSQS()
	seed.IndexerSeed()
}
