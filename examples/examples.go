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
	{"Conway Tribute", ".rle",
		[]string{"#N Conway Tribute",
			"#C A tribute to John Conway on the occasion of his",
			"#C passing",
			"  x = 7, y = 9, rule = b3/s23",
			"2b3o$2bobo$2bobo$3bo$ob3o$bobobo$3bo2bo$2bobo$2bobo!"}},
	{"Sir Robin", ".rle",
		[]string{"#N Sir Robin",
			"#O Adam P. Goucher, Tom Rokicki; 2018",
			"#C The first elementary knightship to be found in Conway's Game of Life.",
			"#C https://conwaylife.com/wiki/Sir_Robin",
			"  x = 31, y = 79, rule = b3/s23",
			"2o$4bo2bo$4bo3bo$6b3o$2b2o6b4o$2bob2o4b4o$bo4bo6b3o$2b4o4b2o3bo$o9b2o$bo3bo$",
			"6b3o2b2o2bo$2b2o7bo4bo$13bob2o$10b2o6bo$11b2ob3obo$10b2o3bo2bo$10bobo2b2o$10bo",
			"2bobobo$10b3o6bo$11bobobo3bo$14b2obobo$11bo6b3o2$11bo9bo$11bo3bo6bo$12bo5b5o$",
			"12b3o$16b2o$13b3o2bo$11bob3obo$10bo3bo2bo$11bo4b2ob3o$13b4obo4b2o$13bob4o4b2o$",
			"19bo$20bo2b2o$20b2o$21b5o$25b2o$19b3o6bo$20bobo3bobo$19bo3bo3bo$19bo3b2o$18bo6b",
			"ob3o$19b2o3bo3b2o$20b4o2bo2bo$22b2o3bo$21bo$21b2obo$20bo$19b5o$19bo4bo$18b3ob3o",
			"$18bob5o$18bo$20bo$16bo4b4o$20b4ob2o$17b3o4bo$24bobo$28bo$24bo2b2o$25b3o$22b2o$",
			"21b3o5bo$24b2o2bobo$21bo2b3obobo$22b2obo2bo$24bobo2b2o$26b2o$22b3o4bo$22b3o4bo$",
			"23b2o3b3o$24b2ob2o$25b2o$25bo2$24b2o$26bo!"}},
	{"Whirlpool", ".rle",
		[]string{"#N Whirlpool",
			"#O Mitch Patenaude mitch@mitchpatenaude.net",
			"#C A long-lived animation of a swirling pattern",
			"#C with a false stagnation.",
			"#C November 30, 2024",
			"  x = 11, y = 11, rule = b3/s23",
			"5b3o$4bo2bo$4bo2bo$6obo$o6b3o$o2bo3bo2bo$b3o6bo$3bob6o$3bo2bo$3bo2bo$3b3o!"}},
	{"Backrake 1", ".rle",
		[]string{"#N Backrake 1",
			"#O Jason Summers",
			"#C An orthogonal period 8 c/2 backrake.",
			"#C www.conwaylife.com/wiki/index.php?title=Backrake_1",
			"x = 27, y = 18, rule = B3/S23",
			"5b3o11b3o5b$4bo3bo9bo3bo4b$3b2o4bo7bo4b2o3b$2bobob2ob2o5b2ob2obobo2b$b",
			"2obo4bob2ob2obo4bob2ob$o4bo3bo2bobo2bo3bo4bo$12bobo12b$2o7b2obobob2o7b",
			"2o$12bobo12b$6b3o9b3o6b$6bo3bo9bo6b$6bobo4b3o11b$12bo2bo4b2o5b$15bo11b",
			"$11bo3bo11b$11bo3bo11b$15bo11b$12bobo!"}},
	{"Blinker Puffer 1", ".rle",
		[]string{"#N Blinker puffer 1",
			"#O Robert Wainwright",
			"#C An orthogonal period 8 c/2 blinker puffer. The first blinker puffer to be found.",
			"#C www.conwaylife.com/wiki/index.php?title=Blinker_puffer_1",
			"x = 9, y = 18, rule = B3/S23",
			"3bo5b$bo3bo3b$o8b$o4bo3b$5o4b4$b2o6b$2ob3o3b$b4o4b$2b2o5b2$5b2o2b$3bo",
			"4bo$2bo6b$2bo5bo$2b6o!"}},
	{"Blinker Ship 1", ".rle",
		[]string{"#N Blinker ship 1",
			"#O Paul Schick",
			"#C A blinker ship created by Paul Schick based on his Schick engine.",
			"#C www.conwaylife.com/wiki/index.php?title=Blinker_ship_1",
			"x = 27, y = 15, rule = B3/S23",
			"10b4o13b$10bo3bo12b$10bo16b$b2o8bo2bo12b$2ob2o22b$b4o3bo18b$2b2o3bob2o",
			"8bo4b3o$6bo3bo8bo4bobo$2b2o3bob2o8bo4b3o$b4o3bo18b$2ob2o22b$b2o8bo2bo",
			"12b$10bo16b$10bo3bo12b$10b4o!"}},
	{"Block Laying Switch Engine", ".rle",
		[]string{"#N Block-laying switch engine",
			"#O Charles Corderman",
			"#C A diagonal period 288 c/12 block-laying puffer.",
			"#C www.conwaylife.com/wiki/index.php?title=Block-laying_switch_engine",
			"x = 29, y = 28, rule = 23/3",
			"18bo10b$b3o8bo5bo10b$o3bo6bo7bo9b$b2o9b4o2b2o9b$3b2ob2o9b3o9b$5b2o11bo",
			"bo8b$19bo7b2o$19bo7b2o11$7b2o20b$7b2o20b7$15b2o12b$15b2o!"}},
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
