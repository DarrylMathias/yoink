package dedup

import (
	"yoink/crawler/validate/hashtable"
	"yoink/utils"
)

func FilterByHash(links []string) ([]string, error){
	var filteredLinks []string
	for _, link := range links{
		// hashurl
		hashedURL, err := utils.HashURL(link)
		if err != nil{
			return nil, err
		}

		// check duplicacy
		seen, err := hashtable.AlreadySeen(hashedURL)
		if err != nil{
			return nil, err
		}
		if(seen){
			continue
		}
		filteredLinks = append(filteredLinks, link)
	}
	return filteredLinks, nil
}