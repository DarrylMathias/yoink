package hashtable

import (
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/redis"

	"gorm.io/gorm"
)

func AlreadySeen(hashedURL string) (bool, error){
	_, err := redis.GetCache(hashedURL)

	// cache hit
	if err == nil{
		return true, nil
	}else{
		// cache miss
		page := new(models.Page)
		db := database.DB
		err := db.Where("url_hash = ?", hashedURL).First(page).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, nil
			}
			return false, err
		}
		err = redis.SetCache(hashedURL, "1")
		if err != nil{
			return false, err
		}
	}

	return true, nil
}