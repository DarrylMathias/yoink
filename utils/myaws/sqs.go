package myaws

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var SqsClient *sqs.Client

func GetSQSClient(){
	sqsClient := sqs.NewFromConfig(*AwsConfig)
	SqsClient = sqsClient
}