package processing

import (
	"yoink/indexer/word_processing/extraction"
	"yoink/indexer/word_processing/tokenizer"
	"yoink/models"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const titleWeight = 5
const descWeight = 3
const wordsWeight = 1

func Process(messages *sqs.ReceiveMessageOutput) ([]models.IndexerOutput, error){
	var output []models.IndexerOutput
	for _, msg := range messages.Messages{
		hash := msg.Body

		// extract text from html
		document, text, err := extraction.ExtractText(hash)
		if err != nil{
			return nil, err
		}

		// tokenize and filter
		title, err := tokenizer.Tokenize(document.Title)
		if err != nil{
			return nil, err
		}
		desc, err := tokenizer.Tokenize(document.Description)
		if err != nil{
			return nil, err
		}
		words, err := tokenizer.Tokenize(text)
		if err != nil{
			return nil, err
		}

		// ewighted hash map
		weightedFreq := make(map[string]int)
		documentLength := 0
		for _, word := range title{
			weightedFreq[word] += 5
			documentLength += 5
		}
		for _, word := range desc{
			weightedFreq[word] += 3
			documentLength += 3
		}
		for _, word := range words{
			weightedFreq[word] += 1
			documentLength += 1
		}

		output = append(output, models.IndexerOutput{
			Hash: *hash,
			WeightedFreq: weightedFreq,
			DocumentLength: documentLength,
		})
	}
	return output, nil
}