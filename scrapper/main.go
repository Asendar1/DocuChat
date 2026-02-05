package main

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if handleDataDir() == false {
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Post("/scrape", func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		urlStr := string(url)

		// after validation, proceed with scraping
		// TODO: add a feature where the user can add multiple URls

		if Scrape(urlStr) == false {
			http.Error(w, "Failed to scrape the URL", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Succeeded"))

	})

	log.Fatal(http.ListenAndServe(":8081", r))
}
