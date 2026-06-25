package store

import (
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

func StoreTF_IDF(indexerOutput []models.IndexerOutput) error{
	
	// reinitializing them to avoid nil pointer dereferencing 
	posting = make(map[string][]models.Posting)
	lexicon = make(map[string]models.Lexicon)

	for _, op := range indexerOutput{
		// document length insertion
		document, err := database.InsertDocLength(&op)
		if err != nil{
			return err
		}

		// insert posting in memory
		for word, freq := range op.WeightedFreq {
			mu.Lock()

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

			mu.Unlock()
		}
		i++

		// check every threshold for new segment which is to pushed to the disk
		threshold, err := strconv.Atoi(env.ConfigValue.PostingThreshold)
		if err != nil{
			return err
		}

		if i >= int64(threshold){
			mu.Lock()

			err := disk.StoreInDisk(&offset, &i, &segmentId, &posting, &lexicon)
			if err != nil{
				return err
			}

			mu.Unlock()
		}

		database.ComputeStatistics(float32(op.DocumentLength))
	}
	return nil
}
