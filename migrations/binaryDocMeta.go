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

var i int64

const noOfPagesInDB = 1_043_092

func task(file *os.File) {
	db := database.DB

	rows, err := db.Model(&models.Page{}).
		Select("id", "document_length").
		Order("id").
		Rows()

	if err != nil {
		fmt.Println("db error:", err)
		return
	}
	defer rows.Close()

	totalRows := 0
	for rows.Next() {
		var page models.Page
		if err := db.ScanRows(rows, &page); err != nil {
			fmt.Println("Scan Error:", err)
			return
		}

		docMeta := models.DocMeta{
			Id:        page.Id,
			DocLength: page.Document_length,
		}

		// Write directly to file
		binary.Write(file, binary.LittleEndian, docMeta)

		totalRows++
		if totalRows%50000 == 0 {
			fmt.Printf("Migrated %d rows...\n", totalRows)
		}
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
