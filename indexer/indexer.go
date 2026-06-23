package indexer

import (
	"fmt"
	"yoink/indexer/store"
	processing "yoink/indexer/word_processing"
	mysqs "yoink/utils/myaws/sqs"
)

func Indexer(sqsURL *string) error{
	// receive message from sqs
	messages, err := mysqs.ReceiveMessage(sqsURL)
	if err != nil{
		return err
	}

	if len(messages.Messages) == 0 {
		return fmt.Errorf("empty sqs queue")
	}
	
	// document processing and words extraction
	indexerOutput, err := processing.Process(messages)
	if err != nil{
		return err
	}

	// db storage
	fmt.Println("indexer output", indexerOutput)
	err = store.StoreTF_IDF(indexerOutput)
	if err != nil{
		return err
	}
	fmt.Println("db push success")

	// delete sqs message
	for _, msg := range messages.Messages{
		if err := mysqs.DeleteMessage(sqsURL, msg); err != nil{
			fmt.Println("delete error:", err)
		}
	}
	return nil
}