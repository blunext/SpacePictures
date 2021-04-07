package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Fetch() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//expected start_date, end_date
		if dateValidator(chi.URLParam(r, "start_date"), chi.URLParam(r, "end_date")) != nil {

		}
	}
}

func dateValidator(start, end string) error {
	return nil
}
