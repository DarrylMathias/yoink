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
	var totalDocLength float32
	var count int
	
	// fetch all hashes
	var hashes []string
	for _, op := range indexerOutput {
		hashes = append(hashes, op.Hash)
	}

	// get id from hash
	pageMap, err := database.GetPageIds(hashes)
	if err != nil {
		return err
	}

	// perform memory posting insertions
	mu.Lock()
	for _, op := range indexerOutput{
		pageId, exists := pageMap[op.Hash]
		if !exists {
			continue
		}
		// insert posting in memory
		for word, freq := range op.WeightedFreq {
			newPosting := models.Posting{
				PageId: pageId,
				TF: int32(freq),
			}
			posting[word] = append(posting[word], newPosting)
		}
		
		i++
		fmt.Printf("doc %d insertion success\n", i)
		
		totalDocLength += float32(op.DocumentLength)
		count++
	}

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

	// update all document lengths in DB
	docLengths := make(map[string]int32)
	for _, op := range indexerOutput {
		docLengths[op.Hash] = int32(op.DocumentLength)
	}
	if err := database.UpdateDocLengths(docLengths); err != nil {
		fmt.Printf("failed to update doc lengths: %v\n", err)
	}

	if count > 0 {
		if err := database.ComputeStatisticsBatch(totalDocLength, count); err != nil {
			fmt.Printf("failed to compute statistics: %v\n", err)
		}
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
