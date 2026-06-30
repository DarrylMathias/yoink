package models

type Data struct{
	Url string `json:"url" gorm:"unique, not null"`
	Title string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Crawl_time int64 `json:"crawl_time"`
	Document_length int32 `json:"document_length"`
	BM25_Rating float64 `json:"BM25_rating"`
}