package golife_test

import (
        "fmt"
        "strings"
        "testing"
        "github.com/pneumaticdeath/golife"
)

var testPattern golife.CellList = golife.CellList{{0, 0}, {1, 0}, {2, 0}}
var testPatternStep golife.CellList = golife.CellList{{1, -1}, {1, 0}, {1, 1}}

func TestPopulationAdd(t *testing.T) {
    pop := make(golife.Population)
    popStep := make(golife.Population)

    if len(pop) != 0 {
        t.Error("New Population isn't empty")
    }

    pop.Add(testPattern)
    if len(pop) != 3 {
        t.Error("Didn't add population properly")
    }

    popStep.Add(testPatternStep)
    if match, errmsg := cmpPops(pop.Step(), popStep); !match {
        t.Error(fmt.Sprintf("Blinker didn't blink properly: %s", errmsg))
    }
}

func cmpPops(pop1, pop2 golife.Population) (bool, string) {
    if len(pop1) != len(pop2) {
        return false, "different lengths"
    }

    minCell1, _ := pop1.BoundingBox()
    minCell2, _ := pop2.BoundingBox()

    for cell, _ := range pop1 {
        translated_cell := golife.Cell{cell.X - minCell1.X + minCell2.X, cell.Y - minCell1.Y + minCell2.Y}
        if !pop2[translated_cell] {
            return false, fmt.Sprintf("cell %v (translated %v) not in pop2", cell, translated_cell)
        }
    }

    return true, ""
}

func TestPopStep(t *testing.T) {
    init_game, err := golife.Load("test_files/turingmachine.rle")
    if err != nil {
        t.Error(err)
        return 
    }

    game_at_1, err := golife.Load("test_files/turingmachine@1.rle")
    if err != nil {
        t.Error(err)
        return
    }

    newPop := init_game.Population.Step()
    if test, errmsg := cmpPops(newPop, game_at_1.Population); !test {
        t.Error(fmt.Sprintf("Didn't calculate next gen properly: %s",errmsg))
    }
}

func TestCellsReader(t *testing.T) {
    cellsReader := strings.NewReader("! foo\n.O.\nO.O\n.O.\n")
    game, err := golife.ReadCells(cellsReader)
    if err != nil {
        t.Error(err)
        return
    }

    if len(game.Comments) != 1 || game.Comments[0] != " foo" {
        t.Error("Failed to parse comment")
    }

    expectedCells := golife.CellList{{1, 0}, {0, 1}, {2, 1}, {1, 2}}
    expectedPop := make(golife.Population)
    expectedPop.Add(expectedCells)

    if matching, errmsg := cmpPops(expectedPop, game.Population); !matching {
        t.Error(fmt.Sprintf("Unexpected game population: %s", errmsg))
    }
}

func BenchmarkGameStep(b *testing.B) {
    game := golife.NewGame()
    game.AddCells(testPattern)

    b.ResetTimer()
    for range b.N {
        game.Next()
    }
}

func BenchmarkBigGame(b *testing.B) {
    game, _ := golife.Load("test_files/turingmachine.rle")

    b.ResetTimer()
    for range b.N {
        game.Next()
    }
}
