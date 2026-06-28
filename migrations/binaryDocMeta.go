package migrations

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"yoink/indexer/store/disk"
	"yoink/models"
	"yoink/utils/database"
	"yoink/utils/env"
)

var i      int64
var mu sync.Mutex

const noOfPagesInDB = 1_043_092

func task(file *os.File) {
	db := database.DB

	for {
		offset := int(atomic.AddInt64(&i, 500) - 500)
		var pages []models.Page
		if offset >= noOfPagesInDB {
			return
		}

		err := db.Limit(500).Offset(offset).Order("id").Find(&pages).Error
		if err != nil {
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
		for _, page := range pages {
			docMeta := models.DocMeta{
				Id:        page.Id,
				DocLength: page.Document_length,
			}
			mu.Lock()
			binary.Write(file, binary.LittleEndian, docMeta)
			mu.Unlock()
		}
	}
}

func MigrateDocMeta() {
	atomic.StoreInt64(&i, 0)
	workers, err := strconv.Atoi(env.ConfigValue.Workers)
	if err != nil {
		panic(err)
	}

	// root path
	root, err := disk.GetRootPath()
	if err != nil {
		panic(err)
	}

	// open docMeta.bin file
	file, err := os.OpenFile(root+"docMeta.bin", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Go(func() { task(file) })
	}
	wg.Wait()
}
