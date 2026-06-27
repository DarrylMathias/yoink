package computation

import (
	"math"
	"yoink/app"
	"yoink/models"

	"github.com/google/uuid"
)

func Compute(token string, tfMapping *map[uuid.UUID]int32, stats models.CorpusStatistics) (*map[uuid.UUID]float64){
	df := len(*tfMapping)
	docs := app.DocumentLengthMap
	// inverse document frequency calcuation
	idf := math.Log(
		( (float64(stats.TotalDocuments) - float64(df) + 0.5) / (float64(df) + 0.5) ) + 1.0,
	)

	// constants
	k1 := 1.6
	b := 0.75

	ranking := make(map[uuid.UUID]float64)
	for pageId, tf := range *tfMapping{
		denom := float64(tf) +
				k1 * (1 - b + b * ( float64(docs[pageId]) / float64(stats.AverageDocLength) ))
		// bm25 final score
		score := idf *
			((float64(tf) * (k1 + 1)) / denom)
		ranking[pageId] = score
	}

	return &ranking
}