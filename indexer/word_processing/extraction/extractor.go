package extraction

import (
	"bytes"
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/myaws/s3"

	"github.com/PuerkitoBio/goquery"
)

func ExtractText(hash *string) (models.Page, string, error){
	// load html from s3
	body, err := s3.GetFile(*hash)
	if err != nil{
		return models.Page{}, "", err
	}

	// load title, description
	document := new(models.Page)
	db := database.DB
	if err := db.Where("url_hash = ?", *hash).Take(document).Error; err != nil{
		return models.Page{}, "", err
	}

	// extract body text
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return models.Page{}, "", err
	}

	text, err := ExtractMeaningfulText(doc)
	return *document, text, err
}