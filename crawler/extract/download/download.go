package download

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"yoink/models"
	"yoink/utils"
)

var (
	ErrNonHTML     = errors.New("content is not html")
	ErrBadStatus   = errors.New("non 200 response")
	ErrTooManyRedirects = errors.New("too many redirects")
)

func DownloadPage(myUrl models.MyURL) ([]byte, error) {
	const MAX_PAGE_SIZE = 10 * 1024 * 1024 // 10 MB

	resp, err := utils.MyGet(myUrl.Url)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	// bad status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrBadStatus, resp.StatusCode)
	}

	// content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("%w: %s", ErrNonHTML, contentType)
	}

	// data
	body, err := io.ReadAll(
		io.LimitReader(resp.Body, MAX_PAGE_SIZE),
	)
	return body, nil
}