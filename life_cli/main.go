package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/pneumaticdeath/golife"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func display(g *golife.Game) {
	min_cell, max_cell := g.Population.BoundingBox()
	width := max_cell.X - min_cell.X + 1

	fmt.Printf("Generation %d\n", g.Generation)
	if max_cell.X >= min_cell.X {
		fmt.Println("Bounding Box", min_cell, "to", max_cell)
		var x, y golife.Coord
		fmt.Printf("+%s+\n", strings.Repeat("-", int(width*2-1)))
		for y = min_cell.Y; y <= max_cell.Y; y++ {
			cell_line := make([]string, width)
			for x = min_cell.X; x <= max_cell.X; x++ {
				if g.Population[golife.Cell{x, y}] {
					cell_line[int(x-min_cell.X)] = "*"
				} else {
					cell_line[int(x-min_cell.X)] = " "
				}
			}
			fmt.Printf("|%s|\n", strings.Join(cell_line, " "))
		}
		fmt.Printf("+%s+\n", strings.Repeat("-", int(width*2-1)))
	}
}

func main() {
	var g *golife.Game

	inputfilePtr := flag.String("input", "", "File to read in")
	outputfilePtr := flag.String("output", "", "File to write output to")
	displayPtr := flag.Bool("display", false, "Display steps")
	generationsPtr := flag.Int("generations", 100, "Number of generations to run")
	pprofPtr := flag.String("pprof", "", "Write profiling output to file")

	flag.Parse()

	if *pprofPtr != "" {
		f, err := os.Create(*pprofPtr)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var err error
	if *inputfilePtr != "" {
		g, err = golife.Load(*inputfilePtr)
		check(err)
	} else {
		cells := []golife.Cell{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {2, 1}}
		g = golife.NewGame()
		g.AddCells(cells)
	}

	for i := 0; i < *generationsPtr && len(g.Population) > 0; i++ {
		if *displayPtr {
			display(g)
		}
		g.Next()
	}
	if *displayPtr {
		display(g)
	}
	if *outputfilePtr != "" {
		err = g.SaveRLE(*outputfilePtr)
		check(err)
	}
}
