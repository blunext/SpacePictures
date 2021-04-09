package linkProvider

import (
	"GogoSpace/app"
	"errors"
	"time"
)

type NasaMock struct {
	dates   map[time.Time]bool
	errDate time.Time
}

func NewNasaMock() *NasaMock {
	n := NasaMock{dates: make(map[time.Time]bool)}
	return &n
}

func (n *NasaMock) GetLink(date time.Time) app.LinkResponse {
	if _, ok := n.dates[date]; ok {
		return app.LinkResponse{Link: date.Format("2006-01-02"), PictureAvailable: true}
	} else if date.Equal(n.errDate) {
		return app.LinkResponse{Err: errors.New("error")}
	}
	return app.LinkResponse{}
}

func (n *NasaMock) AddDate(date time.Time) {
	n.dates[date] = true
}

func (n *NasaMock) RemoveDate(date time.Time) {
	delete(n.dates, date)
}

func (n *NasaMock) SetErrDate(date time.Time) {
	delete(n.dates, date)
	n.errDate = date
}
