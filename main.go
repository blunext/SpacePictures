package main

import (
	"GogoSpace/app"
	"GogoSpace/handler"
	"GogoSpace/linkProvider"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	nasa := linkProvider.NewNasa(apiKey)
	collector := app.NewCollector(nasa, 5)

	r := chi.NewRouter()
	r.Get("/pictures", handler.GetPictures(collector))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("Error ListenAndServe: %v", err)
	}
}
