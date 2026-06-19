package crawler

import (
	"fmt"
	"yoink/crawler/extract"
	"yoink/crawler/store"
	"yoink/crawler/validate"
	mysqs "yoink/utils/myaws/sqs"
)

func Crawl() error {
	// receive message from sqs
	messages, err := mysqs.ReceiveMessage()
	if err != nil{
		return err
	}

	// we defer this so that if any part fails, the message is deleted from sqs always
	defer func(){
		// delete sqs message
		for _, msg := range messages.Messages{
			if err := mysqs.DeleteMessage(msg); err != nil{
				fmt.Println("delete error:", err)
			}
		}
		fmt.Println("deleted sqs message")
	}()

	// phase 1 : normalise and validate messages
	normalizedMessages, err := validate.NormalizeURLData(messages)
	if err != nil{
		return err
	}
	fmt.Println(normalizedMessages)

	// phase 2 : download page and discover new urls
	pages, data, err := extract.ExtractPage(normalizedMessages)
	if err != nil{
		return err
	}

	// phase 3 : s3 and rds storage
	err = store.Store(pages, data)
	if err != nil{
		return err
	}

	return nil
}