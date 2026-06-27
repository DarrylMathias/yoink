package disk

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/user"
	"sort"
	"yoink/indexer/store/database"
	"yoink/models"
)

func Sort(posting *map[string][]models.Posting) []string {
	// extract keys
	var keys []string
	for key := range *posting{
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func StoreInDisk(offset *int64, i *int64, segmentId *int64, posting *map[string][]models.Posting) error{
	// ensure directory exists
	err := os.MkdirAll("/home/ubuntu/indexer_data", 0755)
	if err != nil {
		return err
	}

	// dynamic hostname
	u, err := user.Current()
	if err != nil {
		return err
	}
	host := u.Username

	// Create the posting file
	postingFile, err := os.Create(fmt.Sprintf("/home/%s/indexer_data/posting%d.bin", host, *segmentId))
	if err != nil {
		return err
	}
	defer postingFile.Close()

	// Create the lexicon file
	lexiconFile, err := os.Create(fmt.Sprintf("/home/%s/indexer_data/lexicon%d.bin", host, *segmentId))
	if err != nil {
		return err
	}
	defer lexiconFile.Close()

	// sort according to keys
	keys := Sort(posting)

	// Write the struct to the binary file
	for _, key := range keys{
		post := *posting
		err = binary.Write(postingFile, binary.LittleEndian, post[key])
		if err != nil {
			return err
		}

		// lexicon computation
		length := int64(len(post[key]) * binary.Size(models.Posting{}))
		var bytes [64]byte
		copy(bytes[:], key)
		lex := models.Lexicon{
			Term: bytes,
			Offset: *offset,
			Length: length,
		}
		
		// lexicon saved in file
		err = binary.Write(lexiconFile, binary.LittleEndian, lex)
		if err != nil {
			return err
		} 

		*offset += length
	}

	// re-initiailizations
	*offset = 0
	*i = 0
	fmt.Println("synced to disk, segmenId:", *segmentId)
	database.SetSegmentId(int(*segmentId+1))
	*segmentId++
	
	return nil
}