package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Post("/scrape", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got the msg"))
	})

	log.Fatal(http.ListenAndServe(":8081", r))
}
