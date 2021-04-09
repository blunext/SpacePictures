package handler

import (
	"GogoSpace/app"
	"GogoSpace/linkProvider"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
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
	{"2020-01-01", "2019-12-31", "", false},
	{"2022-01-01", "2022-01-01", "", true},
}

func TestHandler(t *testing.T) {
	nasaMock := linkProvider.NewNasaMock()
	collector := app.NewCollector(nasaMock, 5)

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
		code, _ := makeRequest(t, GetPictures(collector), v.from, v.to)
		assert.Equal(t, v.result, code == http.StatusOK, fmt.Sprintf("handler test failed for dates: %v, %v", v.from, v.to))
	}

	code, _ := makeRequest(t, GetPictures(collector), "", "")
	assert.Equal(t, false, code == http.StatusOK, fmt.Sprintf("handler test failed for empty entries"))

	code, _ = makeRequest(t, GetPictures(collector), "2020-01-01", "")
	assert.Equal(t, false, code == http.StatusOK, fmt.Sprintf("handler test failed bad request, missing end date"))

	code, _ = makeRequest(t, GetPictures(collector), "", "2020-12-31")
	assert.Equal(t, false, code == http.StatusOK, fmt.Sprintf("handler test failed bad request, missing start date"))
}

var concurrencyRanges = []struct {
	from string
	to   string
}{
	{"2000-01-01", "2000-12-31"},
	{"2001-01-01", "2001-12-31"},
	{"2002-01-01", "2002-12-31"},
	{"2003-01-01", "2003-12-31"},
	{"2004-01-01", "2004-12-31"},
	{"2005-01-01", "2005-12-31"},
	{"2006-01-01", "2006-12-31"},
	{"2007-01-01", "2007-12-31"},
	{"2008-01-01", "2008-12-31"},
	{"2009-01-01", "2009-12-31"},
}

func TestConcurrencyPipeline(t *testing.T) {

	// nasa mock returns the same date that was requested instead of the url

	nasaMock := linkProvider.NewNasaMock()
	collector := app.NewCollector(nasaMock, 5)

	for _, v := range concurrencyRanges {
		from, _ := time.Parse("2006-01-02", v.from)
		to, _ := time.Parse("2006-01-02", v.to)

		dates := app.DateRange(from, to)

		for _, date := range dates {
			nasaMock.AddDate(date)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		for _, v := range concurrencyRanges {
			wg.Add(1)
			go concurrentClients(t, v.from, v.to, collector, &wg)
		}
	}

	close := make(chan bool)
	go func() {
		wg.Wait()
		close <- true
	}()

	<-close
}

func concurrentClients(t *testing.T, from string, to string, collector *app.Collector, wg *sync.WaitGroup) {
	_, body := makeRequest(t, GetPictures(collector), from, to)

	responseJson := responseLinks{}
	err := json.Unmarshal([]byte(body), &responseJson)
	if err != nil {
		t.Errorf("error pardsing json: %v", body)
	}

	var dateMap = make(map[time.Time]bool)
	for _, d := range responseJson.Urls {
		date, _ := time.Parse("2006-01-02", d)
		dateMap[date] = true
	}

	fromDate, _ := time.Parse("2006-01-02", from)
	toDate, _ := time.Parse("2006-01-02", to)
	datesRange := app.DateRange(fromDate, toDate)

	for _, d := range datesRange {
		_, ok := dateMap[d]
		assert.True(t, ok, fmt.Sprintf("did not receive date expected %v", d))
	}
	wg.Done()
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
