package store

import (
	"fmt"
	"yoink/models"
	"yoink/utils/myaws/s3"
	"github.com/dustin/go-humanize"
)

func Store(pages []models.Page, data [][]byte) error{
	// store in S3
	for i, page := range pages{
		bytes, err := s3.UploadFile(page.Html_s3_key, data[i])
		if err != nil{
			return err
		}
		fmt.Printf("Uploaded html of url %s of size %s\n", page.Url, humanize.Bytes(uint64(bytes)))
	}

	// store in RDS
	// for i, page := range pages{
	// 	s3.UploadFile(page.Html_s3_key, data[i])
	// 	fmt.Println("Uploaded html of url", page.Url)
	// }

	return nil
}