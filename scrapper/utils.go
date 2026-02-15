package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Asendar1/DocuChat/scrapper/pb"
	"github.com/gocolly/colly/v2"
)

// Scrape takes a URL, scrapes the content (optimized for documentation-style pages like Wikipedia)
// Sends it to the tokenizer microservice for further processing, and returns true if successful.
func Scrape(url string) bool {
	urlHash := sha256.Sum256([]byte(url))
	str := strings.Builder{}

	// TODO add struct and channels to indicate progress and errors (http response writer closes after first write)
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		log.Print("Visiting: " + url)
	})

	success := true

	// Skip navigation, infoboxes, and UI elements
	c.OnHTML("nav, header, footer, aside, .mw-editsection, .navbox, .infobox, .sidebar, .navigation", func(e *colly.HTMLElement) {
		// Ignore these elements
	})

	// Target main content area
	c.OnHTML("article, main, #mw-content-text, [role='main']", func(h *colly.HTMLElement) {
		// Extract headers
		h.ForEach("h1, h2, h3, h4, h5, h6", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			// Skip navigation/metadata headers
			if text != "" && !isMetadataHeader(text) {
				str.WriteString("## " + text + "\n\n")
			}
		})

		// Extract paragraphs
		h.ForEach("p", func(_ int, el *colly.HTMLElement) {
			text := cleanText(el.Text)
			// Filter short or meaningless content
			if text != "" && len(text) > 30 && !isBoilerplate(text) {
				str.WriteString(text + "\n\n")
			}
		})

		// Extract list items (often contain key points)
		h.ForEach("ul > li, ol > li", func(_ int, el *colly.HTMLElement) {
			text := cleanText(el.Text)
			if text != "" && len(text) > 15 {
				str.WriteString("- " + text + "\n")
			}
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %v failed with response: %v\nError: %v", r.Request.URL, r, err)
		success = false
	})

	err := c.Visit(url)
	if err != nil {
		log.Printf("Failed to visit URL: %v", err)
		success = false
	}

	tc, err := NewDocClient("localhost:50051")
	if err != nil {
		log.Printf("Failed to create doc client: %v", err)
		return false
	}
	defer tc.Close()

	fr := &pb.FileRequest{
		Hash: hex.EncodeToString(urlHash[:]),
		Content: str.String(),
	}

	pr, err := tc.client.ProcessFile(context.Background(), fr)
	if err != nil {
		log.Printf("Failed to process file: %v", err)
		return false
	}

	if !PrintResults(pr) {
		return false
	}

	return success
}

func PrintResults(pr *pb.ProcessResponse) bool {
	log.Printf("Success: %v\n", pr.GetSuccess())
	log.Printf("Already Exists: %v\n", pr.GetAlreadyExists())
	log.Printf("Tokens: %v\n", pr.GetTokens())
	log.Printf("Message: %v\n", pr.GetMessage())
	return pr.GetSuccess()
}

func cleanText(text string) string {
	// This is AI generated, however its job is to remove refrence numbers and other boilerplate text.
	citationRegex := regexp.MustCompile(`\[\d+\]|\[citation needed\]|\[edit\]`)
	text = citationRegex.ReplaceAllString(text, "")

	text = strings.Join(strings.Fields(text), " ")

	return strings.TrimSpace(text)
}

func isMetadataHeader(text string) bool {
	// Common Wikipedia metadata sections to skip
	metadataHeaders := []string{
		"References", "External links", "See also", "Notes",
		"Further reading", "Bibliography", "Sources",
		"View source", "View history", "What links here",
		"Related changes", "Upload file", "Permanent link",
		"Page information", "Cite this page", "Get shortened URL",
		"Download QR code", "Download as PDF", "Printable version",
		"Gallery", "Navigation menu",
	}

	textLower := strings.ToLower(text)
	for _, header := range metadataHeaders {
		if strings.ToLower(header) == textLower {
			return true
		}
	}
	return false
}

func isBoilerplate(text string) bool {
	// Skip common boilerplate phrases
	boilerplate := []string{
		"view source",
		"view history",
		"what links here",
		"printable version",
		"permanent link",
	}

	textLower := strings.ToLower(text)
	for _, phrase := range boilerplate {
		if strings.Contains(textLower, phrase) {
			return true
		}
	}
	return false
}

func handleDataDir() bool {
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		log.Printf("Failed to create data directory: %v", err)
		return false
	}
	return true
}
