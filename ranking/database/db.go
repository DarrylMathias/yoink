package database

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"yoink/app"
	"yoink/indexer/store/disk"
	"yoink/models"
	"yoink/utils/database"

	"github.com/google/uuid"
)

func GetCorpusStatistics() (models.CorpusStatistics, error){
	db := database.DB

	stats := new(models.CorpusStatistics)
	err := db.Where("id = 1").First(stats).Error
	if err != nil{
		return models.CorpusStatistics{}, err
	}
	return *stats, nil
}

func GetDocumentLengthBatch() (error) {
	root, err := disk.GetRootPath()
	if err != nil{
		return err
	}

	// open docMeta file
	var docMetas []models.DocMeta
	fp, err := os.Open(root + "docMeta.bin")
	if err != nil{
		return err
	}

	// byte array
	docMetaBytes := make([]byte, binary.Size(models.DocMeta{}))
	var docMeta models.DocMeta
	offset := 0

	for {
		// read one entry
		_, err := fp.ReadAt(docMetaBytes, int64(offset))
		offset += binary.Size(models.DocMeta{})
		if err == io.EOF{
			break
		}
		if err != nil{
			return err
		}

		// read binary data to struct
		err = binary.Read(bytes.NewReader(docMetaBytes), binary.LittleEndian, &docMeta)
		if err != nil{
			return err
		}

		docMetas = append(docMetas, docMeta)
	}

	// convert docMeta array to map
	result := make(map[uuid.UUID]int32)
	for _, docMeta := range docMetas{
		result[docMeta.Id] = docMeta.DocLength
	}	
	app.DocumentLengthMap = result

	return nil
}

func GetDocumentBatch(uuids []uuid.UUID) (*map[uuid.UUID]models.Page, error) {
	db := database.DB
	docMap := make(map[uuid.UUID]models.Page)
	const batchSize = 10000

	for i := 0; i < len(uuids); i += batchSize {
		end := min(i+batchSize, len(uuids))

		var docs []models.Page
		err := db.
			Where("id IN ?", uuids[i:end]).
			Find(&docs).
			Error
		if err != nil {
			return nil, err
		}

		for _, doc := range docs {
			docMap[doc.Id] = doc
		}
	}

	return &docMap, nil
}