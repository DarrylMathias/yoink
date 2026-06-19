package extract

import (
	"fmt"
	"time"
	"yoink/app"
	// "yoink/crawler/extract/dedup"
	"yoink/crawler/extract/download"
	"yoink/crawler/extract/metadata"
	"yoink/models"
	"yoink/utils"
	mysqs "yoink/utils/myaws/sqs"

	"github.com/google/uuid"
)
func ExtractPage(urls []models.MyURL) (pgs []models.Page, data [][]byte, err error){
	var pages []models.Page
	for _, myUrl := range urls{
		// load page html
		byteData, err := download.DownloadPage(myUrl)
		if err != nil{
			return nil, nil, err
		}
		crawlTime := time.Now().Unix()

		// find html title
		title, err := metadata.ExtractTitle(byteData)
		if err != nil{
			return nil, nil, err
		}
		fmt.Println("Title: ", title)

		// find html description
		desc, err := metadata.ExtractDescription(byteData)
		if err != nil{
			return nil, nil, err
		}

		// parse all links
		links, err := metadata.ExtractLinks(byteData)
		if err != nil{
			return nil, nil, err
		}
		fmt.Printf("Parsed %d links\n", len(links))

		// filter links => for now there is no priority, just 30 random links from each page
		const MAX_LINKS_PER_PAGE = 30
		filteredLinks := utils.FilteredURLs(myUrl.Url, links, MAX_LINKS_PER_PAGE)

		// cant afford these many cache and db checks, too expensive, so for now, just pushing to sqs
		// filteredLinks, err = dedup.FilterByHash(filteredLinks)
		// if err != nil{
		// 	return nil, nil, err
		// }
		fmt.Println("Filtered links", filteredLinks)

		// push urls to sqs
		err = mysqs.SendBatchMessage(filteredLinks)
		if err != nil{
			return nil, nil, err
		}
		fmt.Printf("Success sending %d links to SQS\n", len(filteredLinks))
		
		app.Counter += len(filteredLinks)
		
		id, err := uuid.NewRandom()
		if err != nil{
			return nil, nil, err
		}
		pages = append(pages, models.Page{
			Id: id,
			Url: myUrl.Url,
			Url_hash: myUrl.Hash,
			Title: title,
			Description: desc,
			Crawl_time: crawlTime,
			Html_s3_key: myUrl.Hash,
		})
		data = append(data, byteData)
	}
	return pages, data, nil
}