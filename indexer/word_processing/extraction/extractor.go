package extraction

import (
	"bytes"
	"fmt"
	"strings"
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/myaws/s3"

	"github.com/PuerkitoBio/goquery"
)

func ExtractText(hash *string) (string, error){
	// load html from s3
	body, err := s3.GetFile(*hash)
	if err != nil{
		return "", err
	}

	// load title, description
	document := new(models.Page)
	db := database.DB
	if err := db.Where("url_hash = ?", *hash).Take(document).Error; err != nil{
		return "", err
	}

	// extract body text
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	// remove all scripts and styling from client side rendered pages
	doc.Find("script").Remove()
	doc.Find("style").Remove()
	doc.Find("noscript").Remove()
	doc.Find("svg").Remove()

	text := strings.TrimSpace(doc.Find("body").Text())
	fmt.Println(text)
	fmt.Println("title", document.Title)
	fmt.Println("desc", document.Description)

	return text, nil
}