package crawler

import (
	"fmt"
	"sync/atomic"
	"yoink/crawler/extract"
	"yoink/crawler/store"
	"yoink/crawler/validate"
	mysqs "yoink/utils/myaws/sqs"
)

func Crawl(isDiscovering bool, sqsURL *string) error {

	// receive message from sqs
	messages, err := mysqs.ReceiveMessage(sqsURL)
	if err != nil{
		return err
	}

	if len(messages.Messages) == 0 {
		return nil
	}

	// we defer this so that if any part fails, the message is deleted from sqs always
	defer func(){
		// delete sqs message
		for _, msg := range messages.Messages{
			if err := mysqs.DeleteMessage(sqsURL, msg); err != nil{
				fmt.Println("delete error:", err)
			}
		}
		atomic.AddInt64(
			&mysqs.NoOfSQSMessages,
			-int64(len(messages.Messages)),
		)
	}()

	// phase 1 : normalise and validate messages
	normalizedMessages, err := validate.NormalizeURLData(messages)
	if err != nil{
		return err
	}

	// phase 2 : download page and discover new urls
	pages, data, err := extract.ExtractPage(normalizedMessages, isDiscovering, sqsURL)
	if err != nil{
		return err
	}

	fmt.Printf(
		"[counter=%d] processed=%d discovered=%d\n",
		atomic.LoadInt64(&mysqs.NoOfSQSMessages),
		len(normalizedMessages),
		len(pages),
	)

	// phase 3 : s3 and rds storage
	err = store.Store(pages, data)
	if err != nil{
		return err
	}

	return nil
}