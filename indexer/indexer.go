package indexer

import (
	// "fmt"
	"fmt"
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
		fmt.Println("empty sqs queue")
		return nil
	}

	// we defer this so that if any part fails, the message is deleted from sqs always
	// defer func(){
	// 	// delete sqs message
	// 	for _, msg := range messages.Messages{
	// 		if err := mysqs.DeleteMessage(sqsURL, msg); err != nil{
	// 			fmt.Println("delete error:", err)
	// 		}
	// 	}
	// }()
	
	// document processing and words extraction
	weightedHashMaps, documentLengths, err := processing.Process(messages)
	if err != nil{
		return err
	}

	// db storage
	fmt.Println("hashmaps", weightedHashMaps)
	fmt.Println("documentLengths", documentLengths)
	


	return nil
}