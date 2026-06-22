package indexer

import processing "yoink/indexer/word_processing"

func Indexer() error{
	err := processing.Process()
	if err != nil{
		return err
	}
	return nil
}