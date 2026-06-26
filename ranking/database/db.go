package database

import (
	"yoink/models"
	"yoink/utils/database"

	"github.com/google/uuid"
)

func GetCorpusStatistics() (models.CorpusStatistics, error){
	db := database.DB

	stats := new(models.CorpusStatistics)
	err := db.Where("id = 1").First(stats).Error
	if err != nil{
		return models.CorpusStatistics{}, err
	}
	return *stats, nil
}

func GetDocumentBatch(uuids []uuid.UUID) (*map[uuid.UUID]models.Page, error){
	db := database.DB

	var docs []models.Page

	err := db.Where("id IN ?", uuids).Find(&docs).Error
	if err != nil{
		return nil, err
	}

	docMap := make(map[uuid.UUID]models.Page)
	for _, doc := range docs{
		docMap[doc.Id] = doc
	}
	return &docMap, nil
}