package s3

import (
	"bytes"
	"context"
	"io"
	"yoink/utils/env"
	"yoink/utils/myaws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func GetS3Client(){
	s3Client := s3.NewFromConfig(*myaws.AwsConfig)
	S3Client = s3Client
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func UploadFile(key string, body []byte) (int64, error) {
	var bucketName string = env.EnvValue.S3BucketName
	config := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(key + ".html"),
		Body: bytes.NewReader(body),
	}
	_, err := S3Client.PutObject(context.Background(), config)
	if err != nil{
		return 0, err
	}
	return int64(len(body)), nil
}

func GetFile(key string) ([]byte, error){
	var bucketName string = env.EnvValue.S3BucketName
	config := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(key + ".html"),
	}
	output, err := S3Client.GetObject(context.Background(), config)
	if err != nil{
		return nil, err
	}
	body, err := io.ReadAll(output.Body)
	if err != nil{
		return nil, err
	}
	return body, nil
}