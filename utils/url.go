package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/benjaminestes/robots"
	"github.com/temoto/robotstxt"
)

var HTTPClient = &http.Client{
	Timeout: 10*time.Second,
}

func MyGet(url string) (*http.Response, error){
	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		return &http.Response{}, err
	}
	req.Header.Set("User-Agent", "yoinkbot")

	res, err := HTTPClient.Do(req)
	if err != nil{
		return &http.Response{}, err
	}
	return res, nil
}

func IsCrawlable(url string) bool{
	robotsURL, err := robots.Locate(url)
	if err != nil{
		return true
	}
	
	res, err := MyGet(robotsURL)
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

	// only crawl http/https
	if u.Scheme != "" &&
		u.Scheme != "http" &&
		u.Scheme != "https" {
		return "", fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	// lowercase scheme and host
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// remove fragment
	u.Fragment = ""

	// remove default ports
	if strings.HasSuffix(u.Host, ":80") &&
		u.Scheme == "http" {
		u.Host = strings.TrimSuffix(u.Host, ":80")
	}
	if strings.HasSuffix(u.Host, ":443") &&
		u.Scheme == "https" {
		u.Host = strings.TrimSuffix(u.Host, ":443")
	}

	// clean path
	u.Path = path.Clean(u.Path)

	// normalize root paths
	if u.Path == "/" ||
		u.Path == "/." ||
		u.Path == "." {
		u.Path = ""
	}

	// remove tracking params
	q := u.Query()

	trackingParams := []string{
		"utm_source",
		"utm_medium",
		"utm_campaign",
		"utm_content",
		"utm_term",
		"fbclid",
		"gclid",
		"mc_cid",
		"mc_eid",
	}

	for _, param := range trackingParams {
		q.Del(param)
	}

	u.RawQuery = q.Encode()

	// reject URLs without host
	if u.Host == "" {
		return "", fmt.Errorf("empty host")
	}

	// sort query params
	q = u.Query()
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func FilteredURLs(domain string, links []string) []string{
	var filteredLinks []string
	seen := make(map[string]struct{})

	baseURL, err := url.Parse(domain)
	if err != nil {
		return filteredLinks
	}

	normalizedPageURL, err := NormalizeURL(domain)
	if err != nil {
		return filteredLinks
	}

	for _, link := range links {
		if link == "" {
			continue
		}
		linkURL, err := url.Parse(link)
		if err != nil {
			continue
		}

		absoluteURL := baseURL.ResolveReference(linkURL)
		normalizedURL, err := NormalizeURL(
			absoluteURL.String(),
		)
		if err != nil {
			continue
		}

		// skip self-links
		if normalizedURL == normalizedPageURL {
			continue
		}

		// deduplication prevention
		if _, exists := seen[normalizedURL]; exists {
			continue
		}
		seen[normalizedURL] = struct{}{}

		filteredLinks = append(filteredLinks, normalizedURL)
		// fmt.Println(normalizedURL)
	}
	return filteredLinks
}