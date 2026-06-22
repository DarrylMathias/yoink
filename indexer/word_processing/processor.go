package processing

import (
	"fmt"
	"yoink/utils/myaws/s3"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func Process(messages *sqs.ReceiveMessageOutput) error{
	for _, msg := range messages.Messages{
		hash := msg.Body

		// extract text from html in s3
		body, err := s3.GetFile(*hash)
		if err != nil{
			return err
		}

		fmt.Println(string(body))
	}
	return nil
}