package handler

import (
	"GogoSpace/app"
	"GogoSpace/linkProvider"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var ranges = []struct {
	from      string
	to        string
	exclusion string // exclusion date that will not be present in range "from, to"
	result    bool
}{
	{"2020-01-01", "2020-12-31", "", true},
	{"2021-01-01", "2021-12-31", "2021-03-01", false},
}

func TestHandler(t *testing.T) {
	nasaMock := linkProvider.NewNasaMock("key")
	for _, v := range ranges {
		from, _ := time.Parse("2006-01-02", v.from)
		to, _ := time.Parse("2006-01-02", v.to)

		dates := app.DateRange(from, to)

		for _, date := range dates {
			nasaMock.AddDate(date)
		}
		if v.exclusion != "" {
			out, _ := time.Parse("2006-01-02", v.to)
			nasaMock.RemoveDate(out)
		}
		collector := app.NewCollector(nasaMock, 5)
		code, _ := makeRequest(t, GetPictures(collector), v.from, v.to)

		assert.Equal(t, v.result, code == http.StatusOK, fmt.Sprintf("handler test failed for dates: %v, %v", v.from, v.to))
	}

}

func makeRequest(t *testing.T, handler http.Handler, from, to string) (int, string) {
	query := fmt.Sprintf("/pictures?start_date=%s&end_date=%s", from, to)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler.ServeHTTP(r, req)
	return r.Code, r.Body.String()
}
