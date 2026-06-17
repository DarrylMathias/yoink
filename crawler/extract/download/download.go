package download

import (
	"fmt"
	"io"
	"net/http"
	"yoink/utils"
)

func DownloadPageAndStoreInS3(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		return err
	}
	req.Header.Set("User-Agent", "yoinkbot")

	resp, err := utils.HTTPClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}