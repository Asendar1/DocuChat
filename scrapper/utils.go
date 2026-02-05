package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

// TODO: half the file content is empty lines, make files cleaner
func Scrape(url string) bool {
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

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %v failed with response: %v\nError: %v", r.Request.URL, r, err)
		success = false
	})

	err = c.Visit(url)
	if err != nil {
		log.Printf("Failed to visit URL: %v", err)
		success = false
	}

	return success
}

func handleDataDir() bool {
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		log.Printf("Failed to create data directory: %v", err)
		return false
	}
	return true
}
