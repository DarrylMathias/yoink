package migrations

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"yoink/models"
	"yoink/ranking/file"
)

type OldLexicon struct{
	Offset int64 `json:"offset"`
    Length int64 `json:"length"`
}

// script to migrate already indexed lexicons from json to binary

func MigrateToBinaryLexicon() error{
	lexiconFiles, err := file.GetLexiconFiles()
	if err != nil{
		return err
	}

	for _, file := range lexiconFiles{
		if !strings.HasSuffix(file, ".json") {
			continue
		}

		lexicon := new(map[string]OldLexicon)

		lexBytes, err := os.ReadFile(file)
		if err != nil{
			return err
		}

		err = json.Unmarshal(lexBytes, lexicon)
		if err != nil{
			return err
		}

		// Create the lexicon file
		newFile := strings.TrimSuffix(file, ".json") + ".bin"
		lexiconFile, err := os.Create(newFile)
		if err != nil {
			return err
		}
		defer lexiconFile.Close()

		// Sort keys to maintain binary search capability
		var keys []string
		for k := range *lexicon {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, word := range keys {
			oldLexicon := (*lexicon)[word]
			// lexicon computation
			var bytes [64]byte
			copy(bytes[:], word)
			lex := models.Lexicon{
				Term: bytes,
				Offset: oldLexicon.Offset,
				Length: oldLexicon.Length,
			}
			
			// lexicon saved in file
			err = binary.Write(lexiconFile, binary.LittleEndian, lex)
			if err != nil {
				return err
			} 
		}
		fmt.Println("migrated ", file)
	}

	return nil
}