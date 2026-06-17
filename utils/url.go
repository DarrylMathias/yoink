package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/benjaminestes/robots"
	"github.com/temoto/robotstxt"
)

var HTTPClient = &http.Client{
	Timeout: 2*time.Minute,
}

func IsCrawlable(url string) bool{
	robotsURL, err := robots.Locate(url)
	if err != nil{
		return true
	}

	req, err := http.NewRequest("GET", robotsURL, nil)
	if err != nil{
		return true
	}
	req.Header.Set("User-Agent", "yoinkbot")

	res, err := HTTPClient.Do(req)
	if err != nil{
		return true
	}
	defer res.Body.Close()
	
	data, err := robotstxt.FromResponse(res)
	if err != nil{
		return true
	}
	grp := data.FindGroup("yoinkbot")

	return grp.Test(url)
}

func HashURL(url string) (string, error){
	hasher := sha256.New()
	
	// Write the string data to the hasher
	_, err := hasher.Write([]byte(url))
	if err != nil{
		return "", err
	}
	
	// Get the final hash result
	hashedData := hasher.Sum(nil)
	
	// conver to hex
	hexHash := hex.EncodeToString(hashedData)

	return hexHash, nil
}

func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// lowercase scheme and host
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// remove fragment
	u.Fragment = ""

	// remove default ports
	if (u.Scheme == "http" && strings.HasSuffix(u.Host, ":80")) ||
		(u.Scheme == "https" && strings.HasSuffix(u.Host, ":443")) {
		hostParts := strings.Split(u.Host, ":")
		u.Host = hostParts[0]
	}

	// clean path
	u.Path = path.Clean(u.Path)

	// "/" -> ""
	if u.Path == "/" {
		u.Path = ""
	}

	// sort query params
	q := u.Query()
	u.RawQuery = q.Encode()

	return u.String(), nil
}