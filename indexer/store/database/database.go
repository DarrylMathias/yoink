package database

import (
	"fmt"
	"yoink/models"
	"yoink/utils/database"
)

func ComputeStatistics(documentLength float32) error {
	var db = database.DB
	// update stats
	return db.Exec(`
		UPDATE corpus_statistics
		SET
			average_doc_length =
				((average_doc_length * total_documents) + ?)
				/ (total_documents + 1),
			total_documents = total_documents + 1
		WHERE id = 1
	`, documentLength).Error
}

func InsertDocLength(op *models.IndexerOutput) (*models.Page, error) {
		var db = database.DB

		// document table insertion
		document := new(models.Page)
		err := db.Where("url_hash = ?", op.Hash).First(document).Error
		if err != nil{
			return nil, err
		}
		document.Document_length = int32(op.DocumentLength)
		if err := db.Save(document).Error; err != nil {
			return nil, err
		}
		fmt.Println("doc insertion success")
		return document, nil
}