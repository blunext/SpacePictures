package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Get("/pictures", Fetch())

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error ListenAndServe: %v", err)
	}
}
