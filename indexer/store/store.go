package store

import (
	"fmt"
	"strconv"
	"sync"

	"yoink/indexer/store/database"
	"yoink/indexer/store/disk"
	"yoink/models"
	"yoink/utils/env"
)

// global variables for tracking
var offset int64 = 0
var i int64 = 0
var segmentId int64 = 0

// in memory posting and lexicon table
var posting map[string][]models.Posting
var lexicon map[string]models.Lexicon

// mutex for mutual exclusion
var mu sync.Mutex

func Init() {
	// reinitialzing them to avoid nil pointer dereferencing
	posting = make(map[string][]models.Posting)
	lexicon = make(map[string]models.Lexicon)
}

func StoreTF_IDF(indexerOutput []models.IndexerOutput) error{
	for _, op := range indexerOutput{
		// document length insertion
		document, err := database.InsertDocLength(&op)
		if err != nil{
			return err
		}

		// mutex to avoid race conditions
		mu.Lock()
		
		// insert posting in memory
		for word, freq := range op.WeightedFreq {
			newPosting := models.Posting{
				PageId: document.Id,
				TF: int32(freq),
			}
			_, exists := posting[word]
			if !exists{
				posting[word] = []models.Posting{newPosting}
			}else{
				posting[word] = append(posting[word], newPosting)
			}
		}
		
		i++
		fmt.Printf("doc %d insertion success\n", i)

		// check every threshold for new segment which is to pushed to the disk
		threshold, err := strconv.Atoi(env.ConfigValue.PostingThreshold)
		if err != nil{
			mu.Unlock()
			return err
		}

		if i >= int64(threshold){
			err := disk.StoreInDisk(&offset, &i, &segmentId, &posting, &lexicon)
			if err != nil{
				mu.Unlock()
				return err
			}
			// reinitialize maps after successful push to disk
			posting = make(map[string][]models.Posting)
			lexicon = make(map[string]models.Lexicon)
		}

		mu.Unlock()

		database.ComputeStatistics(float32(op.DocumentLength))
	}
	return nil
}

func Flush() error {
	mu.Lock()
	defer mu.Unlock()

	if i > 0 {
		err := disk.StoreInDisk(&offset, &i, &segmentId, &posting, &lexicon)
		if err != nil {
			return err
		}
		// reinitialize maps after successful push to disk
		posting = make(map[string][]models.Posting)
		lexicon = make(map[string]models.Lexicon)
	}
	return nil
}
