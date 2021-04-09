package linkProvider

import (
	"GogoSpace/app"
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

func (n *nasaMock) GetLink(date time.Time) app.LinkResponse {
	if _, ok := n.dates[date]; ok {
		return app.LinkResponse{Link: date.Format("2006-01-02"), PictureAvailable: true}
	}
	return app.LinkResponse{Err: errors.New("error")}
}

func (n *nasaMock) AddDate(date time.Time) {
	n.dates[date] = true
}

func (n *nasaMock) RemoveDate(date time.Time) {
	delete(n.dates, date)
}
