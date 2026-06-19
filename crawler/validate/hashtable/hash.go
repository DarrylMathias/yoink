package hashtable

import (
	"fmt"
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/upstash"

	"gorm.io/gorm"
)

func AlreadySeen(hashedURL string) (bool, error){
	_, err := upstash.GetCache(hashedURL)

	// cache hit
	if err == nil{
		fmt.Println("hash cache hit")
		return true, nil
	}else{
		// cache miss
		fmt.Println("hash cache miss")
		page := new(models.Page)
		db := database.DB
		err := db.Where("url_hash = ?", hashedURL).First(page).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, nil
			}
			return false, err
		}
		err = upstash.SetCache(hashedURL, "1")
		if err != nil{
			return false, err
		}
	}

	return true, nil
}