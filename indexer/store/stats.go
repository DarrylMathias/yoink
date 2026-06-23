package store

import (
	"gorm.io/gorm"
)

func ComputeStatistics(db *gorm.DB, documentLength float32) error {
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