package fetch

import (
	"time"
	"yoink/ranking/file"
	"yoink/ranking/load"

	"github.com/google/uuid"
)

func FetchAllDocs(token string, lexiconTimes *float64, postingTimes *float64) (*map[uuid.UUID]int32, error){
	// doc := make(map[uuid.UUID]float64)

	// get lexicon filenames
	lexiconFiles, err := file.GetLexiconFiles()
	if err != nil{
		return nil, err
	}

	// get lexicon
	t1 := time.Now().UnixMicro()
	lexicon, err := load.LoadOffsets(lexiconFiles, token)
	if err != nil{
		return nil, err
	}
	t2 := time.Now().UnixMicro()
	*lexiconTimes += float64(t2-t1) / 1000.0

	// get postings map
	t1 = time.Now().UnixMicro()
	postings, err := load.LoadPostings(lexicon)
	if err != nil{
		return nil, err
	}
	t2 = time.Now().UnixMicro()
	*postingTimes += float64(t2-t1) / 1000.0

	return postings, nil
}