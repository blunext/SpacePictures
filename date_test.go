package main

import (
	"fmt"
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
	{"2020-01-04", "asas", false},
	{"", "2020-02-05", false},
	{"2020-02-04", "2020-01-04", false},
	{"2020-01-04", "2020-02-04", true},
	{"2020-01-04", "2020-02-05", true},
}

func TestEncodeDecodeMessage(t *testing.T) {
	for _, d := range dates {
		err := dateValidator(d.from, d.to)
		if d.result {
			assert.Nil(t, err, fmt.Sprintf("invalid date test: %s, %s, %v", d.from, d.to, d.result))
		} else {
			assert.Error(t, err, fmt.Sprintf("invalid date test: %s, %s, %v", d.from, d.to, d.result))
		}
	}
}
