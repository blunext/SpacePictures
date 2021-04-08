package app

import (
	"sync"
	"time"
)

type Provider interface {
	GetLink(time.Time) string
}

type request struct {
	response chan string
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

func (c *Collector) ProcessDates(from, to time.Time) []string {
	var wg sync.WaitGroup

	response := make(chan string)

	dates := dateRange(from, to)
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
	for link := range response {
		links = append(links, link)
		//fmt.Println(link)
		wg.Done()
	}
	return links
}

func (c *Collector) process() {
	for {
		req := <-c.requestChannel
		req.response <- c.provider.GetLink(req.date)
	}
}

func dateRange(from, to time.Time) []time.Time {
	var dates []time.Time
	for !from.After(to) {
		dates = append(dates, from)
		from = from.Add(24 * time.Hour)
	}
	return dates
}
