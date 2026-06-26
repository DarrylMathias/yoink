package load

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"yoink/models"
	"yoink/utils/env"

	"github.com/google/uuid"
	"github.com/bcicen/jstream"
)

func LoadOffsets(lexiconFiles []string, word string) (map[string]models.Lexicon, error){
	output := make(map[string]models.Lexicon)

	var wg sync.WaitGroup
	var mu sync.Mutex
	i := 0

	workers, err := strconv.Atoi(env.ConfigValue.Workers)
	if err != nil{
		return nil, err
	}
	for _, file := range lexiconFiles{
		task := func(){
					f, err := os.Open(file)
					if err != nil{
						fmt.Println("error opening file in LoadOffsets, ", err)
						return
					}

					decoder := jstream.NewDecoder(f, 1).EmitKV()
					for mv := range decoder.Stream() {
						if kv, ok := mv.Value.(jstream.KV); ok {
							if kv.Key == word{
								lexiconMap := kv.Value.(map[string]interface{})
								lexicon := models.Lexicon{
									Length: int64(lexiconMap["length"].(float64)),
									Offset: int64(lexiconMap["offset"].(float64)),
								}
								postingFile := strings.ReplaceAll(file, "lexicon", "posting")
								mu.Lock()
								output[strings.ReplaceAll(postingFile, ".json", ".bin")] = lexicon
								mu.Unlock()
							}
						}
					}
				}
		if i < workers{
			wg.Go(task)
			i++
		}else{
			wg.Wait()
			i = 0
		}
	}
	wg.Wait()
	return output, nil
}

func LoadPostings(lexicons map[string]models.Lexicon) (*map[uuid.UUID]int32, error){
	var postings []models.Posting
	for file, lexicon := range lexicons{
		// bytes posting
		postingBytes := make([]byte, lexicon.Length)
		posting := make([]models.Posting, 
			lexicon.Length/int64(binary.Size(models.Posting{})),
		)

		fp, err := os.Open(file)
		if err != nil{
			return nil, err
		}
		defer fp.Close()

		_, err = fp.ReadAt(postingBytes, lexicon.Offset)
		if err != nil{
			return nil, err
		}

		// read binary data to struct
		err = binary.Read(bytes.NewReader(postingBytes), binary.LittleEndian, &posting)
		if err != nil{
			return nil, err
		}

		postings = append(postings, posting...)
	}

	// convert array of posting to map for uniqueness
	mapping := make(map[uuid.UUID]int32)
	for _, posting := range postings{
		mapping[posting.PageId] += posting.TF
	}

	return &mapping, nil
}