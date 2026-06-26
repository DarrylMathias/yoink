package ranking

import (
	"fmt"
	"yoink/indexer/word_processing/tokenizer"
	"yoink/ranking/computation"
	"yoink/ranking/database"
	"yoink/ranking/fetch"

	"yoink/models"

	"github.com/google/uuid"
)

func RankPages(query string) ([]models.Page, error){

	k := 10

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

		// get matching documents in batch by uuid
		var uuids []uuid.UUID
		for pageId := range *tfMapping{
			uuids = append(uuids, pageId)
		}
		if len(uuids) == 0 {
			continue
		}
		documents, err := database.GetDocumentBatch(uuids)
		if err != nil{
			return nil, err
		}
		
		// compute bm25 of each (token, document) pair
	 	bm25MapWord := computation.Compute(token, tfMapping, stats, documents)

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

	return results, nil
}