package mysqs

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
	"yoink/utils/env"
	"yoink/utils/myaws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var SqsClient *sqs.Client
var SQSQueueURL *string
var NoOfSQSMessages int64

func GetSQSClient(){
	sqsClient := sqs.NewFromConfig(*myaws.AwsConfig)
	SqsClient = sqsClient
	if err := GetQueueURL(); err != nil{
		panic(fmt.Errorf("error in fetching queue url --- %s", err.Error()))
	}
	if err := GetNoOfMessages(); err != nil{
		panic(fmt.Errorf("error in getting messages in queue --- %s", err.Error()))
	}
	StartQueueMonitor()
}

func StartQueueMonitor() {
	go func() {
		for {
			if err := GetNoOfMessages(); err != nil {
				fmt.Println("queue monitor:", err)
			}

			time.Sleep(60*time.Second)
		}
	}()
}

func ReceiveMessage() (*sqs.ReceiveMessageOutput, error){
	config := &sqs.ReceiveMessageInput{
		QueueUrl: SQSQueueURL,
		MaxNumberOfMessages: 10,
	}
	message, err := SqsClient.ReceiveMessage(context.Background(), config)
	return message, err
}

func GetNoOfMessages() (error){
	config := &sqs.GetQueueAttributesInput{
        QueueUrl: SQSQueueURL,
        AttributeNames: 
		[]types.QueueAttributeName{
			types.QueueAttributeNameApproximateNumberOfMessages,
		},
    }
    result, err := SqsClient.GetQueueAttributes(context.TODO(), config)
    if err != nil {
        return err
    }

    if approxCount, exists := result.Attributes["ApproximateNumberOfMessages"]; exists {
        val, err := strconv.ParseInt(approxCount, 10, 64)
		atomic.StoreInt64(&NoOfSQSMessages, val)
		if err != nil{
			return err
		}
    } else {
        return fmt.Errorf("Attribute not found.")
    }
	return nil
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

