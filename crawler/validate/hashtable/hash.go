package hashtable

import (
	"fmt"
	"sync/atomic"
	"yoink/app"

	// "yoink/models"
	// "yoink/utils/database"
	"yoink/utils/redis"
	// "gorm.io/gorm"
)

func AlreadySeen(hashedURL string) (bool, error){
	_, err := redis.GetCache(hashedURL)

	// cache hit
	if err == nil{
		atomic.AddInt64(&app.CacheHit, 1)
		return true, nil
	}else{
		// cache miss
		atomic.AddInt64(&app.CacheMiss, 1)

		// db is absolutely hitting limits, so redis will only be the seen set rn

		// page := new(models.Page)
		// db := database.DB
		// err := db.Where("url_hash = ?", hashedURL).First(page).Error

		// if err != nil {
		// 	if err == gorm.ErrRecordNotFound {
		// 		return false, nil
		// 	}
		// 	return false, err
		// }
		err = redis.SetCache(hashedURL, "1")
		if err != nil{
			fmt.Println("redis set failed:", err)
		}
	}

	return true, nil
}