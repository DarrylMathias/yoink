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

	"github.com/google/uuid"
)

var i int64
var mu sync.Mutex

const noOfPagesInDB = 1_043_092

func task(file *os.File) {
	db := database.DB

	lastID := uuid.Nil.String()
	totalRows := 0

	for {
		var pages []models.Page

		err := db.Limit(500).Where("id > ?", lastID).Order("id").Find(&pages).Error
		if err != nil {
			fmt.Println(err)
		}
		if len(pages) == 0 {
			return
		}
		
		for _, page := range pages {
			docMeta := models.DocMeta{
				Id:        page.Id,
				DocLength: page.Document_length,
			}
			mu.Lock()
			binary.Write(file, binary.LittleEndian, docMeta)
			mu.Unlock()
		}

		lastID = pages[len(pages)-1].Id.String()
		totalRows += len(pages)
		fmt.Printf("Migrated %d rows (last ID: %s)\n", totalRows, lastID)
	}
}

func MigrateDocMeta() {
	atomic.StoreInt64(&i, 0)
	fmt.Println("Starting Migration! Configured Workers:", env.ConfigValue.Workers)
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
	file, err := os.OpenFile(root+"docMeta.bin", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
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
