package hashtable

import (
	"fmt"
	"sync/atomic"
	"yoink/app"

	"yoink/models"
	"yoink/utils/database"
	myredis "yoink/utils/redis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AlreadySeen(hashedURL string) (bool, error){
	_, err := myredis.GetCache(hashedURL)

	// cache hit
	if err == nil{
		atomic.AddInt64(&app.CacheHit, 1)
		return true, nil
	}

	// cache miss
	if err == redis.Nil {
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
		if err := myredis.SetCache(hashedURL, "1"); err != nil {
			fmt.Println("redis set failed:", err)
		}

		return false, nil
	}

	// redis failure
	return false, err
}