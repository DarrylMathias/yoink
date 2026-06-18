package s3

import (
	"bytes"
	"context"
	"yoink/utils/env"
	"yoink/utils/myaws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var BucketName string = env.EnvValue.S3BucketName

func GetS3Client(){
	s3Client := s3.NewFromConfig(*myaws.AwsConfig)
	S3Client = s3Client
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func UploadFile(key string, body []byte) (int64, error) {
	config := &s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key: aws.String(key),
		Body: bytes.NewReader(body),
	}
	output, err := S3Client.PutObject(context.Background(), config)
	if err != nil{
		return 0, err
	}
	return (*output.Size), nil
}