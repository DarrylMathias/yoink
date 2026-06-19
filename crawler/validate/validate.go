package validate

import (
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

		// hashurl
		hashedURL, err := utils.HashURL(nURL)
		if err != nil{
			return normalizedURLs, err
		}

		// check duplicacy
		seen, err := hashtable.AlreadySeen(hashedURL)
		if err != nil{
			return normalizedURLs, err
		}
		if(seen){
			continue
		}

		// can crawl
		if(utils.IsCrawlable(nURL)){
			normalizedURLs = append(normalizedURLs, models.MyURL{
				Url: nURL,
				Hash: hashedURL,
			})
		}
	}

	return normalizedURLs, nil
}