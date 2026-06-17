package hashtable

import (
	"yoink/models"
	"yoink/utils/database"

	"gorm.io/gorm"
)

func AlreadySeen(hashedURL string) bool{
	page := new(models.Page)

	db := database.DB
	err := db.Where("url_hash = ?", hashedURL).First(page).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}

	return true
}