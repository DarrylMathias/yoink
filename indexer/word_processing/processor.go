package processing

import (
	"yoink/indexer/word_processing/extraction"
	"yoink/indexer/word_processing/tokenizer"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const titleWeight = 5
const descWeight = 3
const wordsWeight = 1

func Process(messages *sqs.ReceiveMessageOutput) ([]map[string]int, []int, error){
	var weightedFreqSlice []map[string]int
	var documentLengthSlice []int
	for _, msg := range messages.Messages{
		hash := msg.Body

		// extract text from html
		document, text, err := extraction.ExtractText(hash)
		if err != nil{
			return nil, nil, err
		}

		// tokenize and filter
		title, err := tokenizer.Tokenize(document.Title)
		if err != nil{
			return nil, nil, err
		}
		desc, err := tokenizer.Tokenize(document.Description)
		if err != nil{
			return nil, nil, err
		}
		words, err := tokenizer.Tokenize(text)
		if err != nil{
			return nil, nil, err
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

		weightedFreqSlice = append(weightedFreqSlice, weightedFreq)
		documentLengthSlice = append(documentLengthSlice, documentLength)
	}
	return weightedFreqSlice, documentLengthSlice, nil
}