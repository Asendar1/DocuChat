package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Asendar1/NexusProto/DocuChat/gateway/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// File server for static files
	fileServer := http.FileServer(http.Dir("../static/"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	proxies := handlers.NewProxies()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../static/index.html")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/scrape", proxies.HandleScrapeProxy)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
