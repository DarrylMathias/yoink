package ranking

import (
	"slices"
	"sort"
	"yoink/models"
	"yoink/ranking/database"

	"github.com/google/uuid"
)

func Sort(bm25Map *map[uuid.UUID]float64, k int) ([]models.Data, error) {
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
		return nil, nil
	}
	if k > len(ids) {
		k = len(ids)
	}

	// construct unique docSlice
	var result []models.Data
	var uniqueHashes []string
	
	// Fetch documents to check for uniqueness
	for i := 0; i < len(ids) && len(result) < k; {
		end := i + k
		if end > len(ids) {
			end = len(ids)
		}
		
		batchIds := ids[i:end]
		docBatchPtr, err := database.GetDocumentBatch(batchIds)
		if err != nil {
			return nil, err
		}
		docBatch := *docBatchPtr
		
		for _, id := range batchIds {
			doc, exists := docBatch[id]
			if !exists {
				continue
			}
			hash := doc.Url_hash
			if !slices.Contains(uniqueHashes, hash) {
				uniqueHashes = append(uniqueHashes, hash)
				
				res := models.Data{
					Url:             doc.Url,
					Title:           doc.Title,
					Description:     doc.Description,
					Crawl_time:      doc.Crawl_time,
					Document_length: doc.Document_length,
					BM25_Rating:     (*bm25Map)[id],
				}
				result = append(result, res)
				
				if len(result) == k {
					break
				}
			}
		}
		i = end
	}

	return result, nil
}