package golife_test

import (
        "testing"
        "github.com/pneumaticdeath/golife"
)

var testPattern golife.CellList = golife.CellList{{0, 0}, {1, 0}, {2, 0}}

func TestPopulationAdd(t *testing.T) {
    pop := make(golife.Population)

    if len(pop) != 0 {
        t.Error("New Population isn't empty")
    }

    pop.Add(testPattern)
    if len(pop) != 3 {
        t.Error("Didn't add population properly")
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
