package linkProvider

import (
	"GogoSpace/app"
	"time"
)

const nasaUrl = "https://api.nasa.gov/planetary/apod?"

type nasa struct {
	apiKey string
}

func NewNasa(apiKey string) app.Provider {
	n := nasa{apiKey: apiKey}
	return &n
}

func (nasa) GetLink(date time.Time) string {
	return date.Format("2006-01-02")
}
