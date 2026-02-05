package handlers

import (
	"bytes"
	"io"
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
	url, err := io.ReadAll(io.LimitReader(r.Body, 2048))
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("POST", p.scrapeURL, bytes.NewBufferString(string(url)))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach scrapper service", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from scrapper service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
