package models

type IndexerOutput struct{
	Hash string
	WeightedFreq map[string]int
	DocumentLength int
}