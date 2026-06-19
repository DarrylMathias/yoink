package store

import (
	"fmt"
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/myaws/s3"
	"yoink/utils/upstash"

	"gorm.io/gorm/clause"
)

func Store(pages []models.Page, data [][]byte) error{
	// store in S3
	for i, page := range pages{
		_, err := s3.UploadFile(page.Html_s3_key, data[i])
		if err != nil{
			return err
		}
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

		// update redis
		err = upstash.SetCache(page.Url_hash, "1")
		if err != nil{
			fmt.Println("upstash set failed:", err)
		}
	}

	return nil
}