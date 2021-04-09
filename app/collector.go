package app

import (
	"fmt"
	"sync"
	"time"
)

type Provider interface {
	GetLink(time.Time) LinkResponse
}

type LinkResponse struct {
	Link             string
	PictureAvailable bool
	Err              error
}

type request struct {
	response chan LinkResponse
	date     time.Time
}

type Collector struct {
	provider       Provider
	requestChannel chan request
}

func NewCollector(provider Provider, limit int) *Collector {
	c := &Collector{
		provider:       provider,
		requestChannel: make(chan request),
	}
	runGoroutines(c, limit)
	return c
}

func runGoroutines(c *Collector, limit int) {
	for i := 0; i < limit; i++ {
		go func() {
			for {
				req := <-c.requestChannel
				req.response <- c.provider.GetLink(req.date)
			}
		}()
	}
}

func (c *Collector) ProcessDates(from, to time.Time) ([]string, error) {
	var wg sync.WaitGroup

	response := make(chan LinkResponse)

	dates := DateRange(from, to)
	wg.Add(len(dates))

	go func() {
		for _, day := range dates {
			req := request{response: response, date: day}
			c.requestChannel <- req
		}
	}()

	go func() {
		wg.Wait()
		close(response)
	}()

	var links []string
	valid := true
	var lastErr error
	for linkResp := range response {
		// if at least one url is not fetched correctly
		// I assume that all requests are wrong
		if linkResp.Err != nil {
			valid = false
			lastErr = linkResp.Err
		}
		// omitting responses from days where picture is not available
		if linkResp.PictureAvailable {
			links = append(links, linkResp.Link)
		}
		wg.Done()
	}

	if !valid {
		return nil, fmt.Errorf("at least one requested date did not reveive correctly. Maybe you reach OVER_RATE_LIMIT for your api key %w", lastErr)
	}

	return links, nil
}

func DateRange(from, to time.Time) []time.Time {
	var dates []time.Time
	for !from.After(to) {
		dates = append(dates, from)
		from = from.Add(24 * time.Hour)
	}
	return dates
}
