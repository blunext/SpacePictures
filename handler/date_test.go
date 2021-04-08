package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var dates = []struct {
	from   string
	to     string
	result bool
}{
	{"", "", false},
	{"2020-01-04", "", false},
	{"asas", "2020-02-05", false},
	{"2020-01-04", "asas", false},
	{"", "2020-02-05", false},
	{"2020-02-04", "2020-01-04", false},
	{"2020-01-04", "2020-02-04", true},
	{"2020-01-04", "2020-02-05", true},
}

func TestEncodeDecodeMessage(t *testing.T) {
	var date time.Time

	for _, d := range dates {
		from, _, err := dateValidator(d.from, d.to)
		if d.result {
			assert.Nil(t, err, fmt.Sprintf("invalid date test: %s, %s, %v", d.from, d.to, d.result))

			date, err = time.Parse("2006-01-02", d.from)
			assert.Equal(t, date, from, "date conversion failed")
		} else {
			assert.Error(t, err, fmt.Sprintf("invalid date test: %s, %s, %v", d.from, d.to, d.result))
		}

	}
}
