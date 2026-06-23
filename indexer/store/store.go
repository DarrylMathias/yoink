package store

import (
	"fmt"
	"yoink/models"
	"yoink/utils/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// massively optimized it from O(n^2) to O(n) bringing db operations from 4 billion to 4 million for a million corpus
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

		// collect all words
		var words []string
		for key := range op.WeightedFreq{
			words = append(words, key)
		}
		
		// batch check existence of words (basically BATCH GET)
		var foundTerms []models.Term
		err = db.Where("word IN ?", words).Find(&foundTerms).Error
		if err != nil{
			return err
		}
		foundTermsMap := make(map[string]models.Term)
		for _, term := range foundTerms{
			foundTermsMap[term.Word] = term
		}

		// find terms that need to be created then
		var needToBeCreated []models.Term
		for _, word  := range words {
			_, exists := foundTermsMap[word]
			if !exists{
				needToBeCreated = append(needToBeCreated, models.Term{Word: word, DF: 1})
			}
		}
		// BATCH CREATE
		if len(needToBeCreated) > 0 {
			if err := db.CreateInBatches(needToBeCreated, 500).Error; err != nil{
				return err
			}
		}

		// BATCH UPDATE (code for stack overflow)
		values := make([]clause.Expr, 0, len(foundTerms))
		for _, term := range foundTerms {
			term.DF++
			values = append(values, gorm.Expr("(?::bigint, ?::text, ?::integer)", term.Id, term.Word, term.DF))
		}
		valuesExpr := gorm.Expr("?", values)
		valuesExpr.WithoutParentheses = true
		err = db.Exec(
			"UPDATE terms SET word = tmp.word, df = tmp.df FROM (VALUES ?) tmp(id, word, df) WHERE terms.id = tmp.id",
			valuesExpr,
		).Error
		if err != nil{
			return err
		}

		// now postings table
		var allTerms []models.Term
		err = db.Where("word IN ?", words).Find(&allTerms).Error
		if err != nil{
			return err
		}
		termMap := make(map[string]models.Term)
		for _, term := range allTerms{
			termMap[term.Word] = term
		}

		var postings []models.Posting
		for word, value := range op.WeightedFreq{
			term := termMap[word]
			postings = append(postings, models.Posting{
				PageId: document.Id,
				TermId: term.Id,
				TF: int32(value),
			})
		}

		// batch insert postings
		if len(postings) > 0 {
			if err := db.CreateInBatches(postings, 500).Error; err != nil{
				return err
			}
		}

		if err := ComputeStatistics(db, float32(document.Document_length)); err != nil {
			return err
		}
	}
	return nil
}
