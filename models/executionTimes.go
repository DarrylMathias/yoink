package models

type ExecutionTimes struct{
	Tokenize float64 `json:"tokenize"`
	FetchCorpusStats float64 `json:"fetch_corpus_stats"`
	LexiconSeek float64 `json:"lexicon_seek_time"`
	PostingSeek float64 `json:"posting_seek_time"`
	BM25Computation float64 `json:"bm25_computation"`
	Sort float64 `json:"sort"`
}