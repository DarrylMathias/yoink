package models

type Bm25 struct{
	IDF int64
	TF int
	DocLength int
	AvgDocLength int
}