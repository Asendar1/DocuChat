package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Post("/scrape", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got the msg"))
	})

	log.Fatal(http.ListenAndServe(":8081", r))
}
