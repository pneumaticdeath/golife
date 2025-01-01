package examples

import (
	"log"
	"strings"

	"github.com/pneumaticdeath/golife"
)

type Example struct {
	Title    string
	Format   string
	FileData []string
}

var Examples []Example = []Example{
	{"Glider", ".rle",
		[]string{"  x = 3, y = 3, rule = b3/s23", "bo$2bo$3o!"}},
	{"Blinker", ".rle",
		[]string{"  x = 1, y = 3, rule = b3/s23", "o$o$o!"}},
	{"Gosper Glider Gun", ".rle",
		[]string{"#N The Gosper glider gun", "#O Bill Gosper",
			"#C  Discovered by Bill Gosper in November 1970",
			"#   The first known infinitely growing pattern",
			"  x = 36, y = 9, rule = b3/s23",
			"24bo$22bobo$12b2o6b2o12b2o$11bo3bo4b2o12b2o$2o8bo5bo3b2o$2o8bo3bob2o4bobo$10bo",
			"5bo7bo$11bo3bo$12b2o!"}},
}

func LoadExample(e Example) *golife.Game {
	loader := golife.FindReader(e.Format)
	g, err := loader(strings.NewReader(strings.Join(e.FileData, "\n")))
	if err != nil {
		log.Print("Unable to load golife.example:", err)
	} else {
		g.Filename = e.Title // hack, since this isn't really the title of a file
	}

	return g
}
