package mysqs

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
	"yoink/utils/myaws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var SqsClient *sqs.Client
var NoOfSQSMessages int64

func GetSQSClient(){
	sqsClient := sqs.NewFromConfig(*myaws.AwsConfig)
	SqsClient = sqsClient
}

func StartQueueMonitor(queueURL *string) {
	go func() {
		for {
			if err := GetNoOfMessages(queueURL); err != nil {
				fmt.Println("queue monitor:", err)
			}

			time.Sleep(60*time.Second)
		}
	}()
}

func ReceiveMessage(queueURL *string) (*sqs.ReceiveMessageOutput, error){
	config := &sqs.ReceiveMessageInput{
		QueueUrl: queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds: 20,
	}
	message, err := SqsClient.ReceiveMessage(context.Background(), config)
	return message, err
}

func GetNoOfMessages(queueURL *string) (error){
	config := &sqs.GetQueueAttributesInput{
        QueueUrl: queueURL,
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

func SendMessage(queueURL *string, data string) (*sqs.SendMessageOutput, error){
	config := &sqs.SendMessageInput{
		MessageBody: aws.String(data),
		QueueUrl: queueURL,
	}
	output, err := SqsClient.SendMessage(context.Background(), config)
	return output, err
}

func SendBatchMessage(queueURL *string, data []string) (error){
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
			QueueUrl: queueURL,
			Entries: urls,
		}
		_, err := SqsClient.SendMessageBatch(context.Background(), config)
		if err != nil{
			return err
		}
	}
	return nil
}

func DeleteMessage(queueURL *string, input types.Message) error{
	config := &sqs.DeleteMessageInput{
		QueueUrl: queueURL,
		ReceiptHandle: input.ReceiptHandle,
	}
	_, err := SqsClient.DeleteMessage(context.Background(), config)
	return err
}

func GetQueueURL(queueURL string) (*string ,error){
	config := &sqs.GetQueueUrlInput{	
		QueueName: aws.String(queueURL),
	}
	queue, err := SqsClient.GetQueueUrl(context.Background(), config)
	if err != nil{
		return nil, err
	}
	return queue.QueueUrl, nil
}

