package extract

import (
	"fmt"
	"math/rand"
	"yoink/crawler/extract/download"
	"yoink/crawler/extract/metadata"
	"yoink/models"
	"yoink/utils"
	mysqs "yoink/utils/myaws/sqs"
)
func ExtractPage(urls []models.MyURL) error{
	for _, myUrl := range urls{
		// load page html
		data, err := download.DownloadPage(myUrl)
		if err != nil{
			return err
		}

		// find html title
		title, err := metadata.ExtractTitle(data)
		if err != nil{
			return err
		}
		fmt.Println("Title: ", title)

		// find html title
		desc, err := metadata.ExtractDescription(data)
		if err != nil{
			return err
		}
		fmt.Println("Description: ", desc)

		// parse all links
		links, err := metadata.ExtractLinks(data)
		if err != nil{
			return err
		}
		fmt.Printf("Parsed %d links\n", len(links))

		// filter links => for now there is no priority, just 30 random links from each page
		const MAX_LINKS_PER_PAGE = 25

		filteredLinks := utils.FilteredURLs(myUrl.Url, links)
		rand.Shuffle(len(filteredLinks), func(i, j int) {
			filteredLinks[i], filteredLinks[j] = filteredLinks[j], filteredLinks[i]
		})

		filteredLinks = filteredLinks[:MAX_LINKS_PER_PAGE]
		fmt.Println("Filtered links", filteredLinks)

		// push urls to sqs
		for _, link := range filteredLinks{
			_, err := mysqs.SendMessage(link)
			if err != nil{
				return err
			}
		}
		fmt.Printf("Success sending %d links to SQS\n", MAX_LINKS_PER_PAGE)
	}
	return nil
}