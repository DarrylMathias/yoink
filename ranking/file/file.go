package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFiles(extension string) ([]string, error){
	var filesList []string

	// get hostname
	host, err := os.Hostname()
	if err != nil{
		return nil, err
	}
	// get list of all files
	files, err := os.ReadDir(fmt.Sprintf("/home/%s/indexer_data", host))
	if err != nil {
		return nil, err
	}

	// find files list
	for _, file := range files {
		full_path := fmt.Sprintf("/home/%s/indexer_data/%s", host, file.Name())
		if filepath.Ext(full_path) == extension{
			filesList = append(filesList, full_path)
		}
	}
	return filesList, nil
}