package store

import (
	"fmt"
	"yoink/models"
	"yoink/utils/database"

	"gorm.io/gorm"
)

func StoreTF_IDF(indexerOutput []models.IndexerOutput) error{
	var db = database.DB
	for _, op := range indexerOutput{
		// document table insertion
		document := new(models.Page)
		fmt.Println("op", op)
		err := db.Where("url_hash = ?", op.Hash).First(document).Error
		if err != nil{
			return err
		}
		document.Document_length = int32(op.DocumentLength)
		if err := db.Save(document).Error; err != nil {
			return err
		}

		for key, value := range op.WeightedFreq{
			// term table insertion
			term := new(models.Term)
			err := db.Where("word = ?", key).First(term).Error
			if err == gorm.ErrRecordNotFound{
				term = &models.Term{
					Word: key,
					DF: 1,
				}
				if err := db.Create(term).Error; err != nil {
					return err
				}
			}else if err != nil{
				return err
			}else {
				term.DF++
				if err := db.Save(term).Error; err != nil {
					return err
				}
			}
			
			// postings insertion
			posting := new(models.Posting)
			err = db.Where("page_id = ? AND term_id = ?", document.Id, term.Id).First(posting).Error
			if err == gorm.ErrRecordNotFound{
				posting = &models.Posting{
					PageId: document.Id,
					TermId: term.Id,
					TF: int32(value),
				}
				if err := db.Create(posting).Error; err != nil{
					return err
				}
			}else if err != nil{
				return err
			}else{
				posting.TF += int32(value)
				if err := db.Save(posting).Error; err != nil{
					return err
				}
			}
		}
		if err := ComputeStatistics(db, float32(document.Document_length)); err != nil {
			return err
		}
	}
	return nil
}
