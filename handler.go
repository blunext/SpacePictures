package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Fetch() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//expected start_date, end_date
		if dateValidator(r.URL.Query().Get("start_date"), r.URL.Query().Get("end_date")) != nil {

		}
	}
}

func dateValidator(start, end string) error {
	var startTime, endTime time.Time
	var err error

	startTime, err = time.Parse("2006-01-02", start)
	if err != nil {
		log.Printf("cannot parse startTime date: %s", start)
		return fmt.Errorf("cannot parse startTime date")
	}

	endTime, err = time.Parse("2006-01-02", end)
	if err != nil {
		log.Printf("cannot parse endTime date: %s", end)
		return fmt.Errorf("cannot parse endTime date")
	}

	if startTime.After(endTime) {
		log.Printf("startTime date is older than endTime date: %v, %v", startTime, endTime)
		return fmt.Errorf("cannot parse endTime date")
	}
	return nil
}
