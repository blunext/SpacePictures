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

var concurrencyRanges = []struct {
	from      string
	to        string
	exclusion string // exclusion date that will not be present in range "from, to"

}{
	{"2000-01-01", "2000-12-31", "2000-02-01"},
	{"2001-01-01", "2001-12-31", "2001-02-01"},
	{"2002-01-01", "2002-12-31", ""},
	{"2003-01-01", "2003-12-31", ""},
	{"2004-01-01", "2004-12-31", ""},
	{"2005-01-01", "2005-12-31", ""},
	{"2006-01-01", "2006-12-31", ""},
	{"2007-01-01", "2007-12-31", ""},
	{"2008-01-01", "2008-12-31", ""},
	{"2009-01-01", "2009-01-01", "2009-01-01"},
}

func TestStatuses(t *testing.T) {
	nasaMock := preparingData()
	collector := app.NewCollector(nasaMock, 5)

	code, _ := makeRequest(t, GetPictures(collector), "", "")
	assert.Equal(t, false, code == http.StatusOK, fmt.Sprintf("handler test failed for empty entries"))

	code, _ = makeRequest(t, GetPictures(collector), "2020-01-01", "")
	assert.Equal(t, false, code == http.StatusOK, fmt.Sprintf("handler test failed bad request, missing end date"))

	code, _ = makeRequest(t, GetPictures(collector), "", "2020-12-31")
	assert.Equal(t, false, code == http.StatusOK, fmt.Sprintf("handler test failed bad request, missing start date"))

}

func TestConcurrencyPipeline(t *testing.T) {
	var wg sync.WaitGroup

	nasaMock := preparingData()
	collector := app.NewCollector(nasaMock, 5)

	for i := 0; i < 100; i++ {
		for _, v := range concurrencyRanges {
			wg.Add(1)
			concurrentClientTest(t, collector, &wg, v.from, v.to, v.exclusion)
		}
	}

	close := make(chan bool)
	go func() {
		wg.Wait()
		close <- true
	}()

	<-close
}

func preparingData() *linkProvider.NasaMock {
	// nasa mock returns the same date that was requested instead of the url

	nasaMock := linkProvider.NewNasaMock()

	for _, v := range concurrencyRanges {
		from, _ := time.Parse("2006-01-02", v.from)
		to, _ := time.Parse("2006-01-02", v.to)

		dates := app.DateRange(from, to)

		for _, date := range dates {
			nasaMock.AddDate(date)
		}
		exclusion, err := time.Parse("2006-01-02", v.exclusion)
		if err == nil {
			nasaMock.RemoveDate(exclusion)
		}
	}
	return nasaMock
}

func concurrentClientTest(t *testing.T, collector *app.Collector, wg *sync.WaitGroup, from, to, exclusion string) {
	defer wg.Done()

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
	exclusionDate, _ := time.Parse("2006-01-02", exclusion)
	datesRange := app.DateRange(fromDate, toDate)

	for _, d := range datesRange {
		if d.Equal(exclusionDate) {
			continue
		}
		_, ok := dateMap[d]
		assert.True(t, ok, fmt.Sprintf("did not receive date expected %v", d))
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
