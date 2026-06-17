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
		if(hashtable.AlreadySeen(hashedURL)){
			continue
		}
		filteredLinks = append(filteredLinks, link)
	}
	return filteredLinks, nil
}