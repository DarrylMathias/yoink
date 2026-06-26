package file

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

func GetLexiconFiles() ([]string, error){
	var filesList []string

	// get hostname
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	host := u.Username

	// get list of all files
	files, err := os.ReadDir(fmt.Sprintf("/home/%s/indexer_data", host))
	if err != nil {
		return nil, err
	}

	// find files list
	for _, file := range files {
		full_path := fmt.Sprintf("/home/%s/indexer_data/%s", host, file.Name())
		if strings.Contains(full_path, "lexicon"){
			filesList = append(filesList, full_path)
		}
	}
	return filesList, nil
}