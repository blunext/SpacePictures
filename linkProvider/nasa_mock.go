package linkProvider

import (
	"errors"
	"time"
)

type nasaMock struct {
	dates map[time.Time]bool
}

func NewNasaMock() *nasaMock {
	n := nasaMock{dates: make(map[time.Time]bool)}
	return &n
}

func (n *nasaMock) GetLink(date time.Time) (string, error) {
	if _, ok := n.dates[date]; ok {
		return date.Format("2006-01-02"), nil
	}
	return "", errors.New("error")
}

func (n *nasaMock) AddDate(date time.Time) {
	n.dates[date] = true
}

func (n *nasaMock) RemoveDate(date time.Time) {
	delete(n.dates, date)
}
