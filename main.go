package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {
	r := chi.NewRouter()
	r.Get("/pictures", Fetch())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("Error ListenAndServe: %v", err)
	}
}
