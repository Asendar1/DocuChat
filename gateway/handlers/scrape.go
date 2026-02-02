package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type Proxies struct {
	scrapeURL string
	authURL   string
	chatURL   string
}

func NewProxies() *Proxies {
	return &Proxies{
		scrapeURL: "http://localhost:8081/scrape",
		authURL:   "http://localhost:8082/auth",
		chatURL:   "http://localhost:8083/chat",
	}
}

func (p *Proxies) HandleScrapeProxy(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	resp, err := http.Post("http://localhost:8081/scrape", "text/plain", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to contact scrape service", http.StatusInternalServerError)
		log.Print("Check the scrapper service, it may be down")
		return
	}
	defer resp.Body.Close()
}
