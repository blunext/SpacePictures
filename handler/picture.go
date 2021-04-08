package handler

import (
	"GogoSpace/app"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type errorResponse struct {
	Error string `json:"error"`
}

func GetPictures(collector *app.Collector) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		startDate, endDate, err := dateValidator(r.URL.Query().Get("start_date"), r.URL.Query().Get("end_date"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonFromStruct(errorResponse{Error: err.Error()}))
		}

		dates := collector.ProcessDates(startDate, endDate)
		for _, d := range dates {
			fmt.Println(d)
		}

	}
}

func dateValidator(start, end string) (time.Time, time.Time, error) {
	var startTime, endTime time.Time
	var err error

	startTime, err = time.Parse("2006-01-02", start)
	if err != nil {
		log.Printf("cannot parse startTime date: %s", start)
		return startTime, startTime, fmt.Errorf("cannot parse start_date")
	}

	endTime, err = time.Parse("2006-01-02", end)
	if err != nil {
		log.Printf("cannot parse endTime date: %s", end)
		return startTime, endTime, fmt.Errorf("cannot parse end_date")
	}

	if startTime.After(endTime) {
		log.Printf("startTime date is older than endTime date: %v, %v", startTime, endTime)
		return startTime, endTime, fmt.Errorf("cannot parse endTime date")
	}
	return startTime, endTime, nil
}

func jsonFromStruct(s interface{}) []byte {
	j, err := json.Marshal(s)
	if err != nil {

		// TODO: add serverity

		log.Println("cannot marshal %v, err: %v", s, err)
	}
	return j
}
