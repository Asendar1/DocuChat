package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/gocolly/colly/v2"
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

func Scrape(url string) bool {
	if handleDataDir() == false {
		return false
	}
	fileHash := sha256.Sum256([]byte(url))
	fileName := "data/" + hex.EncodeToString(fileHash[:]) + ".txt"
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return false
	}
	defer file.Close()

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		log.Print("Visting: " + url)
	})

	success := true

	// I'll start by getting the whole body for now
	c.OnHTML("body", func(h *colly.HTMLElement) {
		_, err := file.WriteString(h.Text)
		if err != nil {
			log.Printf("Failed to write to file: %v", err)
			success = false
		}
	})

	err = c.Visit(url)
	if err != nil {
		log.Printf("Failed to visit URL: %v", err)
		success = false
	}

	return success
}

func handleDataDir() bool {
	// a better alternative is MkDirAll, however iam a madman
	err := os.Mkdir("data", 0755)
	if err != nil && !os.IsExist(err) {
		log.Printf("Failed to create data directory: %v", err)
		log.Print("Make sure either you have the permission or there is nothing else called \"data\"")
		return false
	}
	return true
}
