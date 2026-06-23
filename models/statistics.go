package models

type CorpusStatistics struct{
	Id   int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	TotalDocuments     uint64 `json:"total_documents"`
    AverageDocLength   float32 `json:"average_documents"`
}