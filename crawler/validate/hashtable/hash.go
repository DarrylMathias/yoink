package hashtable

import (
	"yoink/models"
	"yoink/utils/database"

	"gorm.io/gorm"
)

func AlreadySeen(hashedURL string) (bool, error){
	page := new(models.Page)

	db := database.DB
	err := db.Where("url_hash = ?", hashedURL).First(page).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}