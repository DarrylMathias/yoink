package extract

import (
	"fmt"
	"yoink/crawler/extract/download"
	"yoink/crawler/extract/hashtable"
	"yoink/utils"
	mysqs "yoink/utils/myaws/sqs"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func ExtractURLData() error{
	// receive message from sqs
	messages, err := mysqs.ReceiveMessage()
	if err != nil{
		return err
	}

	// extract html
	for _, msg := range messages.Messages{
		url := aws.ToString(msg.Body)

		// normalize url
		nURL, err := utils.NormalizeURL(url)
		if err != nil{
			return err
		}
		fmt.Printf("Extracting %s\n", url)

		// hashurl
		hashedURL, err := utils.HashURL(url)
		if err != nil{
			return err
		}

		if(hashtable.AlreadySeen(hashedURL)){
			fmt.Println("URL already crawled")
			continue
		}

		// can crawl
		if(utils.IsCrawlable(nURL)){
			fmt.Println("Is crawlable")
			err = download.DownloadPageAndStoreInS3(nURL)
			if err != nil{
				return err
			}
		}else{
			fmt.Printf("Can't crawl page %s :(\n", nURL)
		}
	}

	return nil
}