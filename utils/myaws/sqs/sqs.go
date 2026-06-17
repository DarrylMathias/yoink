package mysqs

import (
	"context"
	"yoink/utils/myaws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var SqsClient *sqs.Client
var SQSQueueURL *string

func GetSQSClient(){
	sqsClient := sqs.NewFromConfig(*myaws.AwsConfig)
	SqsClient = sqsClient
}

func ReceiveMessage() (*sqs.ReceiveMessageOutput, error){
	config := &sqs.ReceiveMessageInput{
		QueueUrl: SQSQueueURL,
	}
	message, err := SqsClient.ReceiveMessage(context.Background(), config)
	return message, err
}

func SendMessage(data string) (*sqs.SendMessageOutput, error){
	config := &sqs.SendMessageInput{
		MessageBody: aws.String(data),
		QueueUrl: SQSQueueURL,
	}
	output, err := SqsClient.SendMessage(context.Background(), config)
	return output, err
}

func GetQueueURL() error{
	config := &sqs.GetQueueUrlInput{	
		QueueName: aws.String("yoink_sqs"),
	}
	queue, err := SqsClient.GetQueueUrl(context.Background(), config)
	if err != nil{
		return err
	}
	SQSQueueURL = queue.QueueUrl
	return nil
}

