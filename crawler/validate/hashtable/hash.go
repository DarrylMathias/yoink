package hashtable

import (
	"fmt"
	"sync/atomic"
	"yoink/app"
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/upstash"

	"gorm.io/gorm"
)

func AlreadySeen(hashedURL string) (bool, error){
	_, err := upstash.GetCache(hashedURL)

	// cache hit
	if err == nil{
		atomic.AddInt64(&app.CacheHit, 1)
		return true, nil
	}else{
		// cache miss
		atomic.AddInt64(&app.CacheMiss, 1)
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
			fmt.Println("upstash set failed:", err)
		}
	}

	return true, nil
}