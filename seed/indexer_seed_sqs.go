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

	"github.com/google/uuid"
)
var i int64

const noOfPagesInDB = 1_043_092

func task(queueURL *string){
	db := database.DB

	lastID := uuid.Nil.String()
	totalRows := 0

	for {
		var pages []models.Page

		err := db.Limit(500).Where("id > ?", lastID).Order("id").Find(&pages).Error   
		if err != nil{
			fmt.Println(err)
		}
		if len(pages) == 0 {
			return
		}

		var msgs []string
		for _, page := range pages{
			msgs = append(msgs, page.Url_hash)
		}
		if err := mysqs.SendBatchMessage(queueURL, msgs); err != nil{
			fmt.Println(err)
		}

		lastID = pages[len(pages)-1].Id.String()
		totalRows += len(pages)
		fmt.Printf("Pushed %d rows (last ID: %s)\n", totalRows, lastID)
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