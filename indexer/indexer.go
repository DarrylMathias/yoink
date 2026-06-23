package indexer

import (
	"errors"
	"fmt"
	"yoink/indexer/store"
	processing "yoink/indexer/word_processing"
	mysqs "yoink/utils/myaws/sqs"
)

var ErrEmptyQueue = errors.New("empty queue")

func Indexer(sqsURL *string) error{
	// receive message from sqs
	messages, err := mysqs.ReceiveMessage(sqsURL)
	if err != nil{
		return err
	}
	if len(messages.Messages) == 0 {
		return ErrEmptyQueue
	}
	
	// document processing and words extraction
	indexerOutput, err := processing.Process(messages)
	if err != nil{
		return err
	}

	// db storage
	err = store.StoreTF_IDF(indexerOutput)
	if err != nil{
		return err
	}

	// delete sqs message
	for _, msg := range messages.Messages{
		if err := mysqs.DeleteMessage(sqsURL, msg); err != nil{
			fmt.Println("delete error:", err)
		}
	}
	return nil
}