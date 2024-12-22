package golife

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "sort"
    "strconv"
    "strings"
    "unicode"
)

const (
    max_line_length = 79
)

type Coord int64

type Cell struct {
    X, Y Coord
}

type Population map[Cell]bool

func check(e error) {
    if e != nil {
        panic(e)
    }
}

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

const neighbor_count_cutoff = 4

func count_neighbors(cell Cell, population Population) int {
    count := 0
    for _, c := range neighbors(cell) {
        if population[c] {
            count += 1
            if count >= neighbor_count_cutoff { // Optimization-- we don't care if it's over 4
                return count
            }
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
                game.History = make([]Population, 0, game.History_size)
            } else {
                game.History = make([]Population, 0, 10)
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
    } else if strings.HasSuffix(filepath, ".life") || strings.HasSuffix(filepath, ".life.txt") {
        return LoadLife(filepath)
    } 
    panic("Unsupported filetype")
}

func LoadRLE(filepath string) *Game {
    var g Game
    g.Init()
    g.Filename = filepath
    g.Comments = make([]string, 0, 10)
    bytes, err := os.ReadFile(filepath)
    check(err)

    contents := string(bytes)
    lines := strings.Split(contents, "\n")

    var count_str strings.Builder
    count := func() int {
        if count_str.Len() == 0 {
            return 1
        } else {
            c, err  := strconv.Atoi(count_str.String())
            check(err)
            count_str.Reset()
            return c
        }
    }

    var x, y, max_x Coord
    var expected_x, expected_y Coord
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
                        expected_x = Coord(ex_x)
                        expected_y = Coord(ex_y)
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
                        y += Coord(count())
                        if max_x < x {
                            max_x = x
                        }
                        x = 0
                    case c == "b":
                        x += Coord(count())
                    case c == "o":
                        num := count()
                        cells := make([]Cell, num)
                        for j := 0; j < num; j++ {
                            var new_cell Cell
                            new_cell.X = x + Coord(j)
                            new_cell.Y = y
                            cells[j] = new_cell
                        }
                        x += Coord(num)
                        g.Population.Add(cells)
                    case c == "!":
                        done = true
                        break
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

    if ! done {
        log.Println("WARN: Did not get terminator at end of RLE file")
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
    
func (population *Population) BoundingBox() (Cell, Cell) {
    var min_cell, max_cell Cell
    min_cell.X = math.MaxInt64
    min_cell.Y = math.MaxInt64
    max_cell.X = math.MinInt64
    max_cell.Y = math.MinInt64

    for cell, present := range *population {
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

type EncodingPair struct {
    symbol string
    count int
}

type CellList []Cell

func (cell_list CellList) Less(i, j int) bool {
    if cell_list[i].Y == cell_list[j].Y {
        return cell_list[i].X < cell_list[j].X
    } else {
        return cell_list[i].Y < cell_list[j].Y
    }
}

func (cell_list CellList) Swap(i, j int) {
    cell_list[i], cell_list[j] = cell_list[j], cell_list[i]
}

func (cell_list CellList) Len() int {
    return len(cell_list)
}

func (game *Game) ExtractRLE() []EncodingPair {
    rle := make([]EncodingPair, 0, 100)

    cells := make(CellList, 0, len(game.Population))
    for c, _ := range game.Population {
        cells = append(cells, c)
    }

    min_cell, _ := game.Population.BoundingBox()

    var last_x Coord = -1
    var last_y Coord = 0
    sort.Sort(cells)
    for i := range cells {
        rel_x := cells[i].X - min_cell.X
        rel_y := cells[i].Y - min_cell.Y

        switch {
        case rel_y > last_y:
            rle = append(rle, EncodingPair{"$", int(rel_y - last_y)})
            last_y = rel_y
            if rel_x > 0 {
                rle = append(rle, EncodingPair{"b", int(rel_x)})
            }
            rle = append(rle, EncodingPair{"o", 1})
        case rel_x == last_x + 1:
            if len(rle) > 0 {
                rle[len(rle)-1].count += 1
            } else {
                rle = append(rle, EncodingPair{"o", 1})
            }
        default:
            rle = append(rle, EncodingPair{"b", int(rel_x - last_x - 1)}, EncodingPair{"o", 1})
        }
        last_x = rel_x
    }
    rle = append(rle, EncodingPair{"!", 1})

    return rle
}

func (game *Game) SaveRLE(filepath string) bool {
    min_cell, max_cell := game.Population.BoundingBox()
    if min_cell.X > max_cell.X {
        return false
    }

    outfile, err := os.Create(filepath)
    check(err)

    defer outfile.Close()

    outwriter := bufio.NewWriter(outfile)
    for i := range game.Comments {
        comment := game.Comments[i]
        if ! strings.HasPrefix(comment, "#") {
            _, err := outwriter.WriteString("#")
            check(err)
        }
        _, err := outwriter.WriteString(comment)
        check(err)
        if ! strings.HasSuffix(comment, "\n") {
            _, err := outwriter.WriteString("\n")
            check(err)
        }
    }
    if len(game.Filename) > 0 {
        _, err := outwriter.WriteString(fmt.Sprintf("#C originally loaded from \"%s\"\n", game.Filename))
        check(err)
    }
    if game.Generation > 0 {
        _, err := outwriter.WriteString(fmt.Sprintf("#C at generation %d\n", game.Generation))
        check(err)
    }
    _, err = outwriter.WriteString(fmt.Sprintf("#C bounded by %d,%d -> %d,%d\n", min_cell.X, min_cell.Y, max_cell.X, max_cell.Y))
    check(err)
    _, err = outwriter.WriteString(fmt.Sprintf("  x = %d, y = %d, rule = b3/s23\n", int(max_cell.X - min_cell.X + 1), int(max_cell.Y - min_cell.Y + 1)))
    check(err)

    var line, newblob strings.Builder

    rle := game.ExtractRLE()
    for i := range rle {
        pair := rle[i]
        if pair.count > 1 {
            newblob.WriteString(fmt.Sprintf("%d%s", pair.count, pair.symbol))
        } else {
            newblob.WriteString(pair.symbol)
        }
        if line.Len() + newblob.Len() > max_line_length {
            line.WriteString("\n")
            _, err := outwriter.WriteString(line.String())
            check(err)
            line.Reset()
        }
        line.WriteString(newblob.String())
        newblob.Reset()
    }
    line.WriteString("\n")
    _, err = outwriter.WriteString(line.String())
    outwriter.Flush()

    return true
}

func LoadLife(filepath string) *Game {
    var game Game
    game.Init()
    game.Filename = filepath
    game.Comments = make([]string, 0, 10)

    f, err := os.ReadFile(filepath)
    check(err)

    content := string(f)
    lines := strings.Split(content, "\n")

    cells := make(CellList, 0, 100)

    for j := range lines {
        line := lines[j]
        chars := strings.Split(line, "")
lineloop:
        for i := range chars {
            char := chars[i]
            switch char {
            case "#":
                game.Comments = append(game.Comments, line[i+1:])
                break lineloop
            case "\n":
            case " ":
            default:
                cells = append(cells, Cell{Coord(i), Coord(j)})
            }
        }
    }

    game.Population.Add(cells)

    return &game
}
