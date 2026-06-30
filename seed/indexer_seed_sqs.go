package seed

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/env"
	mysqs "yoink/utils/myaws/sqs"
	"yoink/utils/resend"
)
var i int64

const noOfPagesInDB = 1_043_092

func task(queueURL *string){
	db := database.DB
	
	rows, err := db.Model(&models.Page{}).
		Select("url_hash").
		Rows()

	if err != nil {
		fmt.Println("DB ERROR:", err)
		return
	}
	defer rows.Close()

	var msgs []string
	totalRows := 0

	for rows.Next() {
		var page models.Page
		if err := database.DB.ScanRows(rows, &page); err != nil {
			fmt.Println("Scan Error:", err)
			return
		}

		msgs = append(msgs, page.Url_hash)

		if len(msgs) == 500 {
			if err := mysqs.SendBatchMessage(queueURL, msgs); err != nil {
				fmt.Println(err)
			}
			totalRows += 500
			msgs = nil // reset slice
			
			if totalRows%50000 == 0 {
				fmt.Printf("Pushed %d rows to SQS...\n", totalRows)
			}
		}
	}

	// Send remaining messages
	if len(msgs) > 0 {
		if err := mysqs.SendBatchMessage(queueURL, msgs); err != nil {
			fmt.Println(err)
		}
		totalRows += len(msgs)
	}
}

func IndexerSeedSQS(){
	atomic.StoreInt64(&i, 0)
	workers, err := strconv.Atoi(env.ConfigValue.Workers)
	if err != nil{
		panic(err)
	}
	queueURL, err := mysqs.GetQueueURL(env.EnvValue.IndexerSqsName)
	if err != nil{
		panic(err)
	}

	var wg sync.WaitGroup
	for w:=0; w<workers; w++{
		wg.Go(func() {task(queueURL)})
	}
	wg.Wait()
	resend.SendEmail("completed push to a million pages sqs", "done indexing")
}