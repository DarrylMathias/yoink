package ranking

import (
	"fmt"
	"time"
	"yoink/indexer/word_processing/tokenizer"
	"yoink/ranking/computation"
	"yoink/ranking/database"
	"yoink/ranking/fetch"

	"yoink/models"

	"github.com/google/uuid"
)


func Init() error {
	err := database.GetDocumentLengthBatch()
	if err != nil{
		return err
	}
	return nil
}

func RankPages(query string, k int) (models.SearchResult, error){
	tTotalStart := time.Now().UnixMilli()
	var result models.SearchResult

	// tokenize query
	t1 := time.Now().UnixMicro()
	tokens, err := tokenizer.Tokenize(query)
	if err != nil{
		return models.SearchResult{}, err
	}
	t2 := time.Now().UnixMicro()
	result.ExecutionTimes.Tokenize = float64(t2-t1) / 1000.0

	// get stats
	t1 = time.Now().UnixMicro()
	stats, err := database.GetCorpusStatistics()
	if err != nil{
		return models.SearchResult{}, err
	}
	t2 = time.Now().UnixMicro()
	result.ExecutionTimes.FetchCorpusStats = float64(t2-t1) / 1000.0

	// whole query map
	bm25Map := make(map[uuid.UUID]float64)
	lexiconTimes := 0.0
	postingTimes := 0.0
	bmTimes := 0.0

	for _, token := range tokens{
		tfMapping, err := fetch.FetchAllDocs(token, &lexiconTimes, &postingTimes)
		if err != nil{
			return models.SearchResult{}, err
		}
		fmt.Println("mapping ", len(*tfMapping))
		
		result.Tokens = append(result.Tokens, map[string]int{token: len(*tfMapping)})

		// compute bm25 of each (token, document) pair
		t1 = time.Now().UnixMicro()
	 	bm25MapWord := computation.Compute(token, tfMapping, stats)
		t2 = time.Now().UnixMicro()
		bmTimes += float64(t2-t1) / 1000.0

		// merge map with central map
		for id, bmRank := range *bm25MapWord{
			bm25Map[id] += bmRank
		}
	}
	result.ExecutionTimes.LexiconSeek = lexiconTimes
	result.ExecutionTimes.PostingSeek = postingTimes
	result.ExecutionTimes.BM25Computation = bmTimes

	// sort and give top k results
	t1 = time.Now().UnixMicro()
	data, err := Sort(&bm25Map, k)
	if err != nil{
		return models.SearchResult{}, err
	}
	result.Data = data
	t2 = time.Now().UnixMicro()
	result.ExecutionTimes.Sort = float64(t2-t1) / 1000.0

	tTotalEnd := time.Now().UnixMilli()
	fmt.Printf("ranking time : %d ms\n", tTotalEnd-tTotalStart)

	return result, nil
}