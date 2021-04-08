package app

import (
	"fmt"
	"sync"
	"time"
)

type Provider interface {
	GetLink(time.Time) (string, error)
}

type linkResponse struct {
	link string
	err  error
}

type request struct {
	response chan linkResponse
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
				link, err := c.provider.GetLink(req.date)
				req.response <- linkResponse{link: link, err: err}
			}
		}()
	}
}

func (c *Collector) ProcessDates(from, to time.Time) ([]string, error) {
	var wg sync.WaitGroup

	response := make(chan linkResponse)

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
	for linkResp := range response {
		// if at least one url is not fetched correctly
		// I assume that all requests are wrong
		if linkResp.err != nil {
			return nil, fmt.Errorf("at least one requested link is wrong: %w", linkResp.err)
		}
		links = append(links, linkResp.link)
		wg.Done()
	}

	return links, nil
}

func (c *Collector) process() {
	for {
		req := <-c.requestChannel
		link, err := c.provider.GetLink(req.date)
		req.response <- linkResponse{link: link, err: err}
	}
}

func DateRange(from, to time.Time) []time.Time {
	var dates []time.Time
	for !from.After(to) {
		dates = append(dates, from)
		from = from.Add(24 * time.Hour)
	}
	return dates
}
