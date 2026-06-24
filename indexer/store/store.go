package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"yoink/models"
	"yoink/utils/database"
)

var offset int64 = 0
var i int64 = 0
var segmentId int64 = 0
var posting map[string][]models.Posting
var lexicon map[string]models.Lexicon
var mu sync.Mutex

func StoreTF_IDF(indexerOutput []models.IndexerOutput) error{
	var db = database.DB

	posting = make(map[string][]models.Posting)
	lexicon = make(map[string]models.Lexicon)

	for _, op := range indexerOutput{
		fmt.Println("op", op)
		// document table insertion
		document := new(models.Page)
		err := db.Where("url_hash = ?", op.Hash).First(document).Error
		if err != nil{
			return err
		}
		document.Document_length = int32(op.DocumentLength)
		if err := db.Save(document).Error; err != nil {
			return err
		}
		fmt.Println("doc insertion success")

		// posting storage in memory
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

		// check 10k for new segment
		if i >= 100{
			mu.Lock()
			// Create the posting file
			file, err := os.Create(fmt.Sprintf("test/posting%d.bin", segmentId))
			if err != nil {
				return err
			}
			defer file.Close()

			// Write the struct to the file
			for word, post := range posting{
				err = binary.Write(file, binary.LittleEndian, post)
				if err != nil {
					return err
				}
				// lexicon computation
				length := int64(len(post) * binary.Size(models.Posting{}))
				lexicon[word] = models.Lexicon{
					Offset: offset,
					Length: length,
				}
				offset += length
			}

			// lexicon stored to lexicon file
			jsonData, err := json.Marshal(lexicon)
			if err != nil {
				return err
			}
			// Write JSON bytes to a file
			err = os.WriteFile(fmt.Sprintf("test/lexicon%d.json", segmentId), jsonData, 0644)
			if err != nil {
				return err
			}

			// re-initiailizations
			offset = 0
			i = 0
			posting = map[string][]models.Posting{}
			lexicon = map[string]models.Lexicon{}
			fmt.Println("disk storage")
			segmentId++
			mu.Unlock()
		}
		ComputeStatistics(db, float32(op.DocumentLength))
	}
	return nil
}
