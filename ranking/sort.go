package ranking

import (
	"sort"
	"yoink/models"
	"yoink/ranking/database"

	"github.com/google/uuid"
)

func Sort(bm25Map *map[uuid.UUID]float64, k int) ([]models.SearchResult, error) {
	bm25MapVal := *bm25Map
	// extract uuids
    var ids []uuid.UUID
    for id := range bm25MapVal {
		ids = append(ids, id)
    }
    
    // sort keys based on their values
    sort.Slice(ids, func(i, j int) bool {
        return bm25MapVal[ids[i]] > bm25MapVal[ids[j]]
    })
    
    // fetch docs of top k matches
	if len(ids) == 0 {
		return []models.SearchResult{}, nil
	}
	if k > len(ids) {
		k = len(ids)
	}

	docMapPtr, err := database.GetDocumentBatch(ids[:k])
	if err != nil{
		return nil, err
	}
	docMap := *docMapPtr

	// construct result array using docMap
	var result []models.SearchResult
	for _, id := range ids[:k]{
		res := models.SearchResult{
			Url: docMap[id].Url,
			Title: docMap[id].Title,
			Description: docMap[id].Description,
			Crawl_time: docMap[id].Crawl_time,
			Document_length: docMap[id].Document_length,
		}
		result = append(result, res)
	}

	return result, nil
}