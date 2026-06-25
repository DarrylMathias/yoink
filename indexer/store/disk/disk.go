package disk

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"yoink/models"
)

func StoreInDisk(offset *int64, i *int64, segmentId *int64, posting *map[string][]models.Posting, lexicon *map[string]models.Lexicon) error{
	// Create the posting file
	file, err := os.Create(fmt.Sprintf("/indexer/posting%d.bin", segmentId))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the struct to the binary file
	for word, post := range *posting{
		err = binary.Write(file, binary.LittleEndian, post)
		if err != nil {
			return err
		}

		// lexicon computation
		length := int64(len(post) * binary.Size(models.Posting{}))
		lex := *lexicon
		lex[word] = models.Lexicon{
			Offset: *offset,
			Length: length,
		}
		*offset += length
	}

	// lexicon stored to lexicon file
	jsonData, err := json.Marshal(lexicon)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("/indexer/lexicon%d.json", segmentId), jsonData, 0644)
	if err != nil {
		return err
	}

	// re-initiailizations
	*offset = 0
	*i = 0
	posting = &map[string][]models.Posting{}
	lexicon = &map[string]models.Lexicon{}
	fmt.Println("synced to disk, segmenId:", *segmentId)
	*segmentId++
	
	return nil
}