package metadata

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func ExtractTitle(data []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(
		bytes.NewReader(data),
	)
	if err != nil {
		return "", err
	}

	title := doc.Find("title").First().Text()
	return title, nil
}

func ExtractDescription(data []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(
		bytes.NewReader(data),
	)
	if err != nil {
		return "", err
	}

	description := doc.Find("meta[name='description']").AttrOr("content", "")

	if description == "" {
		description = doc.Find("meta[property='og:description']").AttrOr("content", "")
	}
	return description, nil
}

func ExtractLinks(data []byte) ([]string, error){
	var parsedURLs []string
	doc, err := goquery.NewDocumentFromReader(
		bytes.NewReader(data),
	)
	if err != nil{
		return nil, err
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			parsedURLs = append(parsedURLs, href)
		}
	})
	return parsedURLs, nil
}