package seed

import (
	"fmt"
	mysqs "yoink/utils/myaws/sqs"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func SeedSQS(){
	seed_urls := []string{"https://en.wikipedia.org", "https://reddit.com", "https://news.ycombinator.com","https://golang.org", "https://darrylmathias.tech"}
	mysqs.GetQueueURL()

	for i, url := range seed_urls{
		output, err := mysqs.SendMessage(url)
		if err != nil{
			panic(fmt.Errorf("error in sqs send message --- %s", err.Error()))
		}
		fmt.Printf("success - url %d: %v\n", i+1, aws.ToString(output.MD5OfMessageBody))
	}
}