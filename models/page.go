package models

import (
	"time"

	"github.com/google/uuid"
)

type Page struct{
	Id uuid.UUID `json:"id" gorm:"primaryKey"`
	Url string `json:"url" gorm:"unique, not null"`
	Url_hash string `json:"url_hash" gorm:"unique, not null"`
	Title string `json:"title" gorm:"not null"`
	Status_code string `json:"status_code"`
	Crawl_time time.Duration `json:"crawl_time"`
	Html_s3_key string `json:"html_s3_key" gorm:"not null"`
}