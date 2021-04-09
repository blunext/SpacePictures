package main

import (
	"GogoSpace/app"
	"GogoSpace/handler"
	"GogoSpace/linkProvider"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	nasa := linkProvider.NewNasa(apiKey)

	limit, err := strconv.Atoi(os.Getenv("CONCURRENT_REQUESTS"))
	if err != nil {
		limit = 5
	}
	collector := app.NewCollector(nasa, limit)

	http.HandleFunc("/pictures", handler.GetPictures(collector))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("Error ListenAndServe: %v", err)
	}
}
