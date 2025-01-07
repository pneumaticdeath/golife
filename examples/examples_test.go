package examples_test

import (
	"github.com/pneumaticdeath/golife/examples"
	"testing"
)

func TestExampleValidity(t *testing.T) {
	var parseCount int
	ex := examples.ListExamples()
	for _, e := range ex {
		if e.Title == "" {
			t.Error("No title")
		}
		g := examples.LoadExample(e)
		if g == nil {
			t.Error("Example ", e.Title, " had nil load")
		} else {
			parseCount += 1
		}
	}
	if parseCount == 0 {
		t.Error("No examples found")
	}
}
