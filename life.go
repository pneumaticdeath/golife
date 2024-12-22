package golife

import (
    "log"
    "math"
    "os"
    "strconv"
    "strings"
    "unicode"
)

type Cell struct {
    X, Y int64
}

type Population map[Cell]bool

func (pop Population) Add(new_cells []Cell) {
    for i := range new_cells {
        pop[new_cells[i]] = true
    }
}

func (current Population) Step() Population {
    nextgen := make(Population)
    look_at := make(Population)
    for cell, present := range current {
        if present {
            look_at[cell] = true
            for _, n := range neighbors(cell) {
                look_at[n] = true
            }
        }
    }

    for cell, present := range look_at {
        if present {
            c := count_neighbors(cell, current)
            if current[cell] && c == 2 || c == 3 {
                nextgen[cell] = true
            }
        }
    }

    return nextgen                
}

func neighbors(cell Cell) []Cell {
    retval := [8]Cell{{cell.X-1, cell.Y-1}, {cell.X, cell.Y-1}, {cell.X+1, cell.Y-1},
                      {cell.X-1, cell.Y}, {cell.X+1, cell.Y}, 
                      {cell.X-1, cell.Y+1}, {cell.X, cell.Y+1}, {cell.X+1, cell.Y+1}}
    return retval[:]
}

func count_neighbors(cell Cell, population Population) int {
    count := 0
    for _, c := range neighbors(cell) {
        if population[c] {
            count += 1
        }
    }
    return count
}

type Game struct {
    Filename string
    Population Population
    History []Population
    History_size int
    Comments []string
    Generation int
}

func (game *Game) Init() {
    game.Population = make(Population)
}

func (game *Game) Next() {
    if game.History_size != 0 {
        if game.History == nil {
            if game.History_size > 0 {
                game.History = make([]Population, game.History_size)
            } else {
                game.History = make([]Population, 10)
            }
        }
        if game.History_size > 0 && len(game.History) >= game.History_size {
            game.History = game.History[-game.History_size + 1:]
        }
        game.History = append(game.History, game.Population)
    } else {
        game.History = nil
    }
    game.Population = game.Population.Step()
    game.Generation += 1
}

func Load(filepath string) *Game {
    if strings.HasSuffix(filepath, ".rle") || strings.HasSuffix(filepath, ".rle.txt") {
        return LoadRLE(filepath)
    } 
    panic("Unsupported filetype")
}

func LoadRLE(filepath string) *Game {
    var g Game
    g.Init()
    g.Filename = filepath
    g.Comments = make([]string, 10)
    bytes, err := os.ReadFile(filepath)
    if err != nil {
        panic(err)
    }

    contents := string(bytes)
    lines := strings.Split(contents, "\n")

    var count_str strings.Builder
    count := func() int {
        if count_str.Len() == 0 {
            return 1
        } else {
            c, err  := strconv.Atoi(count_str.String())
            if err != nil {
                panic(err)
            }
            count_str.Reset()
            return c
        }
    }

    var x, y, max_x int64
    var expected_x, expected_y int64
    done := false

    for line_no := range lines {
        line := lines[line_no]
        if strings.HasPrefix(line, "#") {
            g.Comments = append(g.Comments, strings.TrimPrefix(line, "#"))
        } else {
            if strings.Contains(line, "=") {
                f := func(c rune) bool {
                    return !unicode.IsLetter(c) && !unicode.IsNumber(c)
                }
                fields := strings.FieldsFunc(line, f)
                if len(fields) < 4 || fields[0] != "x" || fields[2] != "y" {
                    log.Printf("ERROR: Unable to parse line: %s", line)
                } else {
                    ex_x, err1 := strconv.Atoi(fields[1])
                    ex_y, err2 := strconv.Atoi(fields[3])
                    if err1 != nil || err2 != nil {
                        log.Printf("ERROR: Unable to parse x and/or y: %s", line)
                    } else {
                        expected_x = int64(ex_x)
                        expected_y = int64(ex_y)
                    }
                }
            } else {
                chars := strings.Split(strings.TrimSpace(line), "")
                for i := range chars {
                    c := chars[i]
                    switch {
                    case unicode.IsNumber(rune(c[0])):
                        count_str.WriteString(c)
                    case c == "$":
                        y += int64(count())
                        if max_x < x {
                            max_x = x
                        }
                        x = 0
                    case c == "b":
                        x += int64(count())
                    case c == "o":
                        num := count()
                        cells := make([]Cell, num)
                        for j := 0; j < num; j++ {
                            var new_cell Cell
                            new_cell.X = x + int64(j)
                            new_cell.Y = y
                            cells[j] = new_cell
                        }
                        x += int64(num)
                        g.Population.Add(cells)
                    case c == "!":
                        done = true
                    default:
                        log.Printf("ERROR: Unknown code point %s", c)
                    }
                }
            }
        }
        if done {
            break
        }
    }

    if max_x < x {
        max_x = x
    }
    y += 1
    if max_x != expected_x || y != expected_y {
        log.Printf("WARN: Expected board of %dx%d got %dx%d", expected_x, expected_y, max_x, y)
    }
    return &g
}
    
func (game Game) BoundingBox() (Cell, Cell) {
    var min_cell, max_cell Cell
    min_cell.X = math.MaxInt64
    min_cell.Y = math.MaxInt64
    max_cell.X = math.MinInt64
    max_cell.Y = math.MinInt64

    for cell, present := range game.Population {
        if present {
            if cell.X < min_cell.X {
                min_cell.X = cell.X
            }
            if cell.Y < min_cell.Y {
                min_cell.Y = cell.Y
            }
            if cell.X > max_cell.X {
                max_cell.X = cell.X
            }
            if cell.Y > max_cell.Y {
                max_cell.Y = cell.Y
            }
        }
    }

    return min_cell, max_cell
}
