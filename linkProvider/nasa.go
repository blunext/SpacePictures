package linkProvider

import (
	"SpacePictures/app"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// GET https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY&date=2019-12-06

const (
	nasaUrl         = "https://api.nasa.gov/planetary/apod?"
	rateLimitHeader = "X-RateLimit-Remaining"
)

type Response struct {
	Url string `json:"url"`
}

type nasa struct {
	apiKey string
}

func NewNasa(apiKey string) app.Provider {
	n := nasa{apiKey: apiKey}
	return &n
}

func (n nasa) GetLink(date time.Time) app.LinkResponse {

	url := n.getUrl(date)

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return app.LinkResponse{Err: fmt.Errorf("cannot get url %v: %w", url, err)}
	}

	switch {
	case resp.StatusCode == http.StatusNotFound:
		// there are some days that picture is not present
		// https://api.nasa.gov/planetary/apod?api_key=DEMO_API&date=2020-06-10
		return app.LinkResponse{}
	case resp.StatusCode != http.StatusOK:
		// Note: X-RateLimit-Remaining header seems not to show proper value...
		return app.LinkResponse{Err: fmt.Errorf("wrong status %d, on request %s, X-RateLimit-Remaining %v",
			resp.StatusCode, url, resp.Header.Get(rateLimitHeader))}
	}

	body, err := io.ReadAll(resp.Body)

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return app.LinkResponse{Err: fmt.Errorf("marshal response error: %w", err)}
	}

	return app.LinkResponse{Link: response.Url, PictureAvailable: true}
}

func (n nasa) getUrl(date time.Time) string {
	params := url.Values{}
	params.Set("api_key", n.apiKey)
	params.Set("date", date.Format("2006-01-02"))

	return nasaUrl + params.Encode()
}
