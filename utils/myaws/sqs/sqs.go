package mysqs

import (
	"context"
	"strconv"
	"yoink/utils/env"
	"yoink/utils/myaws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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

func SendBatchMessage(data []string) (error){
	// chunking since aws allows only 10 msgs per batch
	for i:=0; i<len(data); i+=10{
		end := i + 10
		if end > len(data) {
			end = len(data)
		}

		var urls []types.SendMessageBatchRequestEntry
		for j, url := range data[i:end] {
			urls = append(urls, types.SendMessageBatchRequestEntry{
				Id: aws.String(strconv.Itoa(j)),
				MessageBody: aws.String(url),
			})
		}
		config := &sqs.SendMessageBatchInput{
			QueueUrl: SQSQueueURL,
			Entries: urls,
		}
		_, err := SqsClient.SendMessageBatch(context.Background(), config)
		if err != nil{
			return err
		}
	}
	return nil
}

func DeleteMessage(input types.Message) error{
	config := &sqs.DeleteMessageInput{
		QueueUrl: SQSQueueURL,
		ReceiptHandle: input.ReceiptHandle,
	}
	_, err := SqsClient.DeleteMessage(context.Background(), config)
	return err
}

func GetQueueURL() error{
	config := &sqs.GetQueueUrlInput{	
		QueueName: aws.String(env.EnvValue.SqsName),
	}
	queue, err := SqsClient.GetQueueUrl(context.Background(), config)
	if err != nil{
		return err
	}
	SQSQueueURL = queue.QueueUrl
	return nil
}

