package main

import (
	"net/http"
)

func ValidateURL(url string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil || req.URL.Scheme == "" || req.URL.Host == "" {
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ValidateURL/1.0)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false
	}
	return true
}
