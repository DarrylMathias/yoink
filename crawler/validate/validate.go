package validate

import (
	"fmt"
	"yoink/crawler/validate/hashtable"
	"yoink/models"
	"yoink/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NormalizeURLData(messages *sqs.ReceiveMessageOutput) ([]models.MyURL, error){
	normalizedURLs := []models.MyURL{}
	for _, msg := range messages.Messages{
		url := aws.ToString(msg.Body)

		// normalize url
		nURL, err := utils.NormalizeURL(url)
		if err != nil{
			return normalizedURLs, err
		}
		fmt.Printf("Validating %s\n", url)

		// hashurl
		hashedURL, err := utils.HashURL(url)
		if err != nil{
			return normalizedURLs, err
		}

		// check duplicacy
		if(hashtable.AlreadySeen(hashedURL)){
			fmt.Println("URL already crawled")
			continue
		}

		// can crawl
		if(utils.IsCrawlable(nURL)){
			fmt.Println("Is crawlable")
			normalizedURLs = append(normalizedURLs, models.MyURL{
				Url: nURL,
				Hash: hashedURL,
			})
		}else{
			fmt.Printf("Can't crawl page %s :(\n", nURL)
		}
	}

	return normalizedURLs, nil
}