package logging

import (
	"fmt"
	"sync/atomic"
	"time"
	mysqs "yoink/utils/myaws/sqs"
	"yoink/utils/resend"
)

func StartHeartbeat(hit *int64, miss *int64) {
	go func() {
		for {
			fmt.Printf(
				"[HEARTBEAT] frontier=%d hits=%d misses=%d\n",
				atomic.LoadInt64(&mysqs.NoOfSQSMessages),
				atomic.LoadInt64(hit),
				atomic.LoadInt64(miss),
			)

			time.Sleep(time.Minute)
		}
	}()
}

func SendHearbeatMail(hit *int64, miss *int64) {
	go func() {
		for {
			err, mailId := resend.SendEmail(
				fmt.Sprintf(
					"[HEARTBEAT] frontier=%d hits=%d misses=%d\n",
					atomic.LoadInt64(&mysqs.NoOfSQSMessages),
					atomic.LoadInt64(hit),
					atomic.LoadInt64(miss),
				),
				"Crawling updates",
			)

			if err != nil {
				fmt.Printf("heartbeat mail failed: %v\n", err)
			} else {
				fmt.Printf("heartbeat mail sent: %s\n", mailId)
			}

			time.Sleep(6 * time.Hour)
		}
	}()
}