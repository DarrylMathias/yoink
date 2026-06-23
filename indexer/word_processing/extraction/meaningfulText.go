package extraction

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ai generated function cause i didnt have the energy to handle all these string formatting
func ExtractMeaningfulText(doc *goquery.Document) (string, error){
	// remove all scripts and styling from client side rendered pages
	doc.Find("script").Remove()
	doc.Find("style").Remove()
	doc.Find("noscript").Remove()
	doc.Find("svg").Remove()

	// remove common layout noise
	doc.Find("header").Remove()
	doc.Find("footer").Remove()
	doc.Find("nav").Remove()
	doc.Find("aside").Remove()

	var builder strings.Builder

	// prioritize content-bearing elements
	doc.Find("h1,h2,h3,h4,h5,h6,p,article,section,li,blockquote,pre").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		builder.WriteString(text)
		builder.WriteString("\n")
	})

	text := builder.String()

	// fallback if nothing useful found
	if strings.TrimSpace(text) == "" {
		text = doc.Find("body").Text()
	}

	// normalize whitespace
	text = strings.Join(strings.Fields(text), " ")

	return text, nil
}