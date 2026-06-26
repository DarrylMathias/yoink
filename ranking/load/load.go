package load

import (
	"bytes"
	"encoding/binary"
	"os"
	"yoink/models"

	"github.com/google/uuid"
)

func LoadOffsets(lexiconFiles []string, word string) (map[string]models.Lexicon, error){
	output := make(map[string]models.Lexicon)
	lexicons := make(map[string]models.Lexicon)

	for _, file := range lexiconFiles{
		
		// find total no of bytes in file
		fi, err := os.Stat(file)
		if err != nil {
			return nil, err
		}
		noOfLexiconBytes := fi.Size()

		// create a file pointer
		fp, err := os.Open(file)
		if err != nil{
			return nil, err
		}
		defer fp.Close()

		// get the size of one lexicon struct
		lexiconSize := binary.Size(models.Lexicon{})
		l := int64(0)
		r := noOfLexiconBytes/int64(lexiconSize)

		lexBytes := make([]byte, lexiconSize)
		lexicon := new(models.Lexicon)

		// perform binary search on data
		for l <= r{
			mid := l + (r-l)/2
			fp.ReadAt(lexBytes, mid*int64(lexiconSize))

			// read binary data to struct
			err = binary.Read(bytes.NewReader(lexBytes), binary.LittleEndian, lexicon)
			if err != nil{
				return nil, err
			}

			// conditions
			if string(lexicon.Term[:]) == word{
				break
			}
			if string(lexicon.Term[:]) < word{
				r = mid - 1
			}else{
				l = mid + 1
			}
		}
		if l <= r{
			lexicons[word] = *lexicon
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