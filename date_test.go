package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var dates = []struct {
	from   string
	to     string
	result bool
}{
	{"", "", false},
	{"2020-01-04", "", false},
	{"asas", "2020-02-05", false},
	{"2020-01-04", "asas", true},
	{"", "2020-02-05", false},
	{"2020-01-05", "2020-02-04", false},
	{"2020-01-04", "2020-02-04", true},
	{"2020-01-04", "2020-02-05", true},
}

func TestEncodeDecodeMessage(t *testing.T) {
	for _, d := range dates {
		err := dateValidator(d.from, d.to)
		if d.result {
			assert.Nil(t, err, "invalid date test")
		} else {
			assert.Error(t, err, "invalid date test")
		}
	}
}
