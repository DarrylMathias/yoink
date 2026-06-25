package indexer

import (
	"errors"
	"fmt"
	"yoink/indexer/store"
	processing "yoink/indexer/word_processing"
	mysqs "yoink/utils/myaws/sqs"
)

var ErrEmptyQueue = errors.New("empty queue")

func Indexer(sqsURL *string) (int, error) {
	// receive message from sqs
	messages, err := mysqs.ReceiveMessage(sqsURL)
	if err != nil{
		return 0, err
	}
	if len(messages.Messages) == 0 {
		return 0, ErrEmptyQueue
	}
	
	// document processing and words extraction
	indexerOutput, err := processing.Process(messages)
	if err != nil{
		return 0, err
	}

	// indexer storage
	err = store.StoreTF_IDF(indexerOutput)
	if err != nil{
		return 0, err
	}

	// delete sqs messages
	if len(messages.Messages) > 0 {
		if err := mysqs.DeleteBatchMessages(sqsURL, messages.Messages); err != nil{
			fmt.Println("delete batch error:", err)
		} else {
			fmt.Printf("deleted %d sqs messages\n", len(messages.Messages))
		}
	}
	return len(messages.Messages), nil
}