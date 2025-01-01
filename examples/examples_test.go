package examples_test

import (
        "testing"
        "github.com/pneumaticdeath/golife/examples"
)

func TestExampleValidity(t *testing.T) {
	for index := range examples.Examples{
		e := examples.Examples[index]
		if e.Title == "" {
			t.Error("No title at index", index)
		}
		g := examples.LoadExample(e)
		if g == nil {
			t.Error("Example",e.Title,"had nil load")
		}
	}
}
