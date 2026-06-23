package processing

import (
	"yoink/indexer/word_processing/extraction"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func Process(messages *sqs.ReceiveMessageOutput) error{
	for _, msg := range messages.Messages{
		hash := msg.Body
		extraction.ExtractText(hash)
		
	}
	return nil
}