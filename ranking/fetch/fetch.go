package fetch

import (
	"yoink/ranking/file"
	"yoink/ranking/load"

	"github.com/google/uuid"
)

func FetchAllDocs(token string) (*map[uuid.UUID]int32, error){
	// doc := make(map[uuid.UUID]float64)

	// get lexicon filenames
	lexiconFiles, err := file.GetFiles(".json")
	if err != nil{
		return nil, err
	}

	// get lexicon
	lexicon, err := load.LoadOffsets(lexiconFiles, token)
	if err != nil{
		return nil, err
	}

	// get postings map
	postings, err := load.LoadPostings(lexicon)
	if err != nil{
		return nil, err
	}

	return postings, nil
}