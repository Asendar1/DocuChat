package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func createProxy(target string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("Proxy error: %v", e)
		http.Error(w, "Proxy error", http.StatusBadGateway)
	}

	return proxy, nil
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	scrapeProxy, err := createProxy("http://localhost:8081")
	if err != nil {
		log.Fatalf("Failed to create scrape proxy: %v", err)
	}

	// File server for static files
	fileServer := http.FileServer(http.Dir("../static/"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../static/index.html")
	})

	r.HandleFunc("/api/v1/scrape", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/scrape"
		scrapeProxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
