package load

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"strings"
	"yoink/models"

	"github.com/google/uuid"
)

func LoadOffsets(lexiconFiles []string, word string) (map[string]models.Lexicon, error){
	output := make(map[string]models.Lexicon)

	for _, file := range lexiconFiles{
		lex := new(map[string]models.Lexicon)

		// lexicon in bytes
		bytes, err := os.ReadFile(file)
		if err != nil{
			return nil, err
		}
		// lexicon in json
		err = json.Unmarshal(bytes, lex)
		if err != nil{
			return nil, err
		}
		lexicons := *lex
		lexicon, exists := lexicons[word]
		if exists{
			postingFile := strings.ReplaceAll(file, "lexicon", "posting")
			output[strings.ReplaceAll(postingFile, ".json", ".bin")] = lexicon
		}
	}
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