package linkProvider

import (
	"GogoSpace/app"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// GET https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY&date=2019-12-06

const nasaUrl = "https://api.nasa.gov/planetary/apod?"

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

func (n nasa) GetLink(date time.Time) (string, error) {

	url := n.getUrl(date)

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("cannot get url %v: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wrong status: %d, on request %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("marshal response error: %w", err)
	}

	return response.Url, nil
}

func (n nasa) getUrl(date time.Time) string {
	params := url.Values{}
	params.Set("api_key", n.apiKey)
	params.Set("date", date.Format("2006-01-02"))

	return nasaUrl + params.Encode()
}
