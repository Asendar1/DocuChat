package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/scrape", http.StatusMovedPermanently)
	})

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

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		log.Printf("Scraper server is running on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down scrapper server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown scrapper server: %v", err)
	}

	log.Println("Scrapper server exiting")

}
