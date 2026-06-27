package disk

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/user"
	"sort"
	"yoink/indexer/store/database"
	"yoink/models"

	"github.com/google/uuid"
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

func GetRootPath() (string, error){
	// dynamic hostname
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	host := u.Username
	return fmt.Sprintf("/home/%s/indexer_data/", host), nil
}

func StoreDiskLength(indexerOutput []models.IndexerOutput, pageMap map[string]uuid.UUID) error{
	// root path
	root, err := GetRootPath()
	if err != nil{
		return err
	}

	// open docMeta.bin file
	file, err := os.OpenFile(root + "docMeta.bin", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// append entry to file
	for _, op := range indexerOutput{
		docMeta := models.DocMeta{
			Id: pageMap[op.Hash],
			DocLength: int32(op.DocumentLength),
		}
		binary.Write(file, binary.LittleEndian, docMeta)
	}
	return nil
}

func StoreInDisk(offset *int64, i *int64, segmentId *int64, posting *map[string][]models.Posting) error{
	// ensure directory exists
	err := os.MkdirAll("/home/ubuntu/indexer_data", 0755)
	if err != nil {
		return err
	}

	// get root path
	root, err := GetRootPath()
	if err != nil{
		return err
	}

	// Create the posting file
	postingFile, err := os.Create(fmt.Sprintf("%sposting%d.bin", root, *segmentId))
	if err != nil {
		return err
	}
	defer postingFile.Close()

	// Create the lexicon file
	lexiconFile, err := os.Create(fmt.Sprintf("%slexicon%d.bin", root, *segmentId))
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