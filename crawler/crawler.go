package crawler

import (
	"fmt"
	"sync/atomic"
	"yoink/crawler/extract"
	"yoink/crawler/store"
	"yoink/crawler/validate"
	"yoink/utils/env"
	mysqs "yoink/utils/myaws/sqs"

	"github.com/aws/aws-sdk-go-v2/aws"
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
		// delete sqs messages
		if len(messages.Messages) > 0 {
			if err := mysqs.DeleteBatchMessages(sqsURL, messages.Messages); err != nil{
				fmt.Println("delete batch error:", err)
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

	// final phase => push to indexer sqs queue
	var msgs []string
	for _, msg := range normalizedMessages{
		msgs = append(msgs, msg.Hash)
	}
	mysqs.SendBatchMessage(aws.String(env.EnvValue.IndexerSqsName), msgs)

	return nil
}