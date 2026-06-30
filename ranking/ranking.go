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

func RankPages(query string, k int) ([]models.Page, error){
	t1 := time.Now().UnixMilli()

	// tokenize query
	tokens, err := tokenizer.Tokenize(query)
	if err != nil{
		return nil, err
	}

	// get stats
	stats, err := database.GetCorpusStatistics()
	if err != nil{
		return nil, err
	}

	// whole query map
	bm25Map := make(map[uuid.UUID]float64)

	for _, token := range tokens{
		tfMapping, err := fetch.FetchAllDocs(token)
		if err != nil{
			return nil, err
		}
		fmt.Println("mapping ", len(*tfMapping))
		
		// compute bm25 of each (token, document) pair
	 	bm25MapWord := computation.Compute(token, tfMapping, stats)

		// merge map with central map
		for id, bmRank := range *bm25MapWord{
			bm25Map[id] += bmRank
		}
	}

	// sort and give top k results
	results, err := Sort(&bm25Map, k)
	if err != nil{
		return nil, err
	}

	t2 := time.Now().UnixMilli()
	fmt.Printf("ranking time : %d ms\n", t2-t1)

	return results, nil
}