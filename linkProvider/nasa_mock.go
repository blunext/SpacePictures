package linkProvider

import (
	"errors"
	"time"
)

type nasaMock struct {
	apiKey string
	dates  map[time.Time]bool
}

func NewNasaMock(apiKey string) *nasaMock {
	n := nasaMock{apiKey: apiKey, dates: make(map[time.Time]bool)}
	return &n
}

func (n *nasaMock) GetLink(date time.Time) (string, error) {
	if _, ok := n.dates[date]; ok {
		return "ok", nil
	}
	return "", errors.New("error")
}

func (n *nasaMock) AddDate(date time.Time) {
	n.dates[date] = true
}

func (n *nasaMock) RemoveDate(date time.Time) {
	delete(n.dates, date)
}
