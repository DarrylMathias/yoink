package seed

import (
	"context"
	"fmt"
	"yoink/app"
	mysqs "yoink/utils/myaws/sqs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SeedSQS(){
	app.App()
	seed_urls := []string{"https://en.wikipedia.org", "https://reddit.com", "https://news.ycombinator.com","https://golang.org"}

	sqsClient := mysqs.SqsClient
	mysqs.GetQueueURL()

	for i, url := range seed_urls{
		config := &sqs.SendMessageInput{
			MessageBody: aws.String(url),
			QueueUrl: mysqs.SQSQueueURL,
		}
		output, err := sqsClient.SendMessage(context.Background(), config)
		if err != nil{
			panic(fmt.Errorf("error in sqs send message --- %s", err.Error()))
		}
		fmt.Printf("success - url %d: %v\n", i+1, aws.ToString(output.MD5OfMessageBody))
	}
}