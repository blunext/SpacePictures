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

type responseLinks struct {
	Urls []string `json:"urls"`
}

func GetPictures(collector *app.Collector) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		startDate, endDate, err := dateValidator(r.URL.Query().Get("start_date"), r.URL.Query().Get("end_date"))
		if err != nil {
			failedResponse(w, http.StatusBadRequest, err)
			return
		}

		var links []string
		links, err = collector.ProcessDates(startDate, endDate)
		if err != nil {
			failedResponse(w, http.StatusNotFound, err)
			return
		}

		resp := responseLinks{Urls: links}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonFromStruct(resp))
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
		return startTime, endTime, fmt.Errorf("startTime date is older than endTime date: %v, %v", startTime, endTime)
	}
	return startTime, endTime, nil
}

func jsonFromStruct(s interface{}) []byte {
	j, err := json.Marshal(s)
	if err != nil {
		log.Printf("cannot marshal %v, err: %v\n", s, err)
		return []byte(fmt.Sprintf("{\"error\": \"cannot marshal: %v, %v\"}", s, err))
	}
	return j
}

func failedResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write(jsonFromStruct(errorResponse{Error: err.Error()}))
}
