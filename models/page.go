package models

import (
	"github.com/google/uuid"
)

type Page struct{
	Id uuid.UUID `json:"id" gorm:"primaryKey"`
	Url string `json:"url" gorm:"unique, not null"`
	Url_hash string `json:"url_hash" gorm:"unique, not null"`
	Title string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Crawl_time int64 `json:"crawl_time"`
	Html_s3_key string `json:"html_s3_key" gorm:"not null"`
	Document_length int32 `json:"document_length"`
}