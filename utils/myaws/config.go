package myaws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var AwsConfig *aws.Config

func GetConfig() error{
	if AwsConfig == nil{
		awsConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-south-1"))
		if err != nil{
			return err
		}
		AwsConfig = &awsConfig
	}
	return nil
}
