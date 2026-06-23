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

	for {
		offset := int(atomic.AddInt64(&i, 500) - 500)
		var pages []models.Page
		if offset >= noOfPagesInDB {
			return
		}

		err := db.Limit(500).Offset(offset).Find(&pages).Error   
		if err != nil{
			fmt.Println(err)
		}
		if len(pages) == 0 {
			return
		}
		fmt.Println(
			"offset:",
			offset,
			"rows:",
			len(pages),
		)

		var msgs []string
		for _, page := range pages{
			msgs = append(msgs, page.Url_hash)
		}
		if err := mysqs.SendBatchMessage(queueURL, msgs); err != nil{
			fmt.Println(err)
		}
	}
}

func IndexerSeed(){
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