package store

import (
	"fmt"
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/myaws/s3"
	"yoink/utils/upstash"

	"github.com/dustin/go-humanize"
	"gorm.io/gorm/clause"
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
	for _, page := range pages{
		db := database.DB
		err := db.Clauses(
			clause.OnConflict{
				DoNothing: true,
			},
		).Create(page).Error
		if err != nil{
			return err
		}
		fmt.Println("Metadata stored in RDS successfully")

		// update redis
		err = upstash.SetCache(page.Url_hash, "1")
		if err != nil{
			return err
		}
		fmt.Println("upstash redis updated successfully")
	}

	return nil
}