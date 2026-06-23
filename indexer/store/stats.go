package store

import (
	"yoink/models"

	"gorm.io/gorm"
)

func ComputeStatistics(db *gorm.DB, documentLength float32) error{
	stats := new(models.CorpusStatistics)
	err := db.Where("id = 1").First(stats).Error
	if err != nil{
		return err
	}
	avgDocLength := (stats.AverageDocLength*float32(stats.TotalDocuments) + documentLength)/float32(stats.TotalDocuments+1)
	stats.AverageDocLength = avgDocLength
	stats.TotalDocuments++

	if err := db.Save(stats).Error; err != nil{
		return err
	}
	return nil
}