package fetch

import (
	"fmt"
	"time"
	"yoink/ranking/file"
	"yoink/ranking/load"

	"github.com/google/uuid"
)

func FetchAllDocs(token string) (*map[uuid.UUID]int32, error){
	// doc := make(map[uuid.UUID]float64)

	// get lexicon filenames
	t1 := time.Now().UnixMilli()
	lexiconFiles, err := file.GetLexiconFiles()
	if err != nil{
		return nil, err
	}
	t2 := time.Now().UnixMilli()
	fmt.Printf("get lexicon filenames took %d ms\n", t2-t1)

	// get lexicon
	t1 = time.Now().UnixMilli()
	lexicon, err := load.LoadOffsets(lexiconFiles, token)
	if err != nil{
		return nil, err
	}
	t2 = time.Now().UnixMilli()
	fmt.Printf("load offsets took %d ms\n", t2-t1)

	// get postings map
	t1 = time.Now().UnixMilli()
	postings, err := load.LoadPostings(lexicon)
	if err != nil{
		return nil, err
	}
	t2 = time.Now().UnixMilli()
	fmt.Printf("load postings took %d ms\n", t2-t1)

	return postings, nil
}