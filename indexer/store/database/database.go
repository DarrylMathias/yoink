package database

import (
	"yoink/models"
	"yoink/utils/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ComputeStatisticsBatch(totalDocumentLength float32, count int) error {
	var db = database.DB
	// made compute statistics batch now
	return db.Exec(`
		UPDATE corpus_statistics
		SET
			average_doc_length =
				((average_doc_length * total_documents) + ?)
				/ (total_documents + ?),
			total_documents = total_documents + ?
		WHERE id = 1
	`, totalDocumentLength, count, count).Error
}

func GetSegmentId() (int, error){
	var db = database.DB
	stats := new(models.CorpusStatistics)

	err := db.Where("id = 1").First(&stats).Error
	if err != nil{
		return 0, err
	}

	return int(stats.SegmentId), nil
}

func SetSegmentId(segmentId int) error{
	var db = database.DB

	err := db.Model(&models.CorpusStatistics{}).Where("id = ?", 1).Update("segment_id", segmentId).Error
	if err != nil{
		return err
	}
	return nil
}

func GetPageIds(hashes []string) (map[string]uuid.UUID, error) {
	var pages []models.Page
	db := database.DB
	
	err := db.Select("id, url_hash").Where("url_hash IN ?", hashes).Find(&pages).Error
	if err != nil {
		return nil, err
	}
	pageMap := make(map[string]uuid.UUID)
	for _, p := range pages {
		pageMap[p.Url_hash] = p.Id
	}
	return pageMap, nil
}

func UpdateDocLengths(docLengths map[string]int32) error {
	var db = database.DB
	return db.Transaction(func(tx *gorm.DB) error {
		for hash, length := range docLengths {
			if err := tx.Model(&models.Page{}).Where("url_hash = ?", hash).Update("document_length", length).Error; err != nil {
				return err
			}
		}
		return nil
	})
}