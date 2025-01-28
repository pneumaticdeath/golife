package golife

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/mitchellh/copystructure"
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

func (pop Population) Size() int {
	return len(pop)
}

func (current Population) Step() Population {
	nextgen := make(Population)
	neighbor_count := make(map[Cell]int8)
	for cell, present := range current {
		if present {
			for _, n := range neighbors(cell) {
				neighbor_count[n] += 1
			}
		}
	}

	for cell, count := range neighbor_count {
		if count == 3 || count == 2 && current[cell] {
			nextgen[cell] = true
		}
	}

	return nextgen
}

func neighbors(cell Cell) []Cell {
	retval := [8]Cell{{cell.X - 1, cell.Y - 1}, {cell.X, cell.Y - 1}, {cell.X + 1, cell.Y - 1},
		{cell.X - 1, cell.Y}, {cell.X + 1, cell.Y},
		{cell.X - 1, cell.Y + 1}, {cell.X, cell.Y + 1}, {cell.X + 1, cell.Y + 1}}
	return retval[:]
}

type Game struct {
	Filename    string
	Population  Population
	History     []Population
	HistorySize int
	Name        string
	Author      string
	Comments    []string
	Generation  int
}

func NewGame() *Game {
	var game Game
	game.Init()
	return &game
}

func (game *Game) Size() int {
	return game.Population.Size()
}

func (game *Game) Copy() *Game {
	copied, copyerr := copystructure.Copy(game)
	if copyerr != nil {
		log.Fatal("Unable to copy game")
		return nil
	}
	newgame, casterr := copied.(Game)
	if casterr {
		log.Fatal("Unable to cast the copy")
		return nil
	}
	return &newgame
}

func (game *Game) Init() {
	game.Population = make(Population)
	game.Comments = make([]string, 0, 10)
	game.History = make([]Population, 0, 10)
}

func (game *Game) SetHistorySize(size int) {
	if game.History == nil {
		if size > 0 {
			game.History = make([]Population, 0, size)
		} else if size == 0 {
			game.History = nil
		} else {
			game.History = make([]Population, 0, 10)
		}
	} else if size > 0 && len(game.History) > size {
		game.History = game.History[len(game.History)-size:]
	}

	game.HistorySize = size
}

func (game *Game) AddCell(cell Cell) {
	game.Population[cell] = true
}

func (game *Game) AddCells(cells []Cell) {
	game.Population.Add(cells)
}

func (game *Game) RemoveCell(cell Cell) {
	delete(game.Population, cell)
}

func (game *Game) HasCell(cell Cell) bool {
	return game.Population[cell]
}

func (game *Game) Next() {
	if game.HistorySize != 0 {
		if game.History == nil {
			if game.HistorySize > 0 {
				game.History = make([]Population, 0, game.HistorySize)
			} else {
				game.History = make([]Population, 0, 10)
			}
		}
		if game.HistorySize > 0 && len(game.History) >= game.HistorySize {
			game.History = game.History[len(game.History)-game.HistorySize+1:]
		}
		game.History = append(game.History, game.Population)
	} else {
		game.History = nil
	}
	game.Population = game.Population.Step()
	game.Generation += 1
}

func (game *Game) Previous() error {
	if game.HistorySize == 0 || len(game.History) == 0 {
		return io.EOF
	}

	prevPop := game.History[len(game.History)-1]
	game.History = game.History[:len(game.History)-1]
	game.Population = prevPop
	game.Generation -= 1
	return nil
}

func Load(filepath string) (*Game, error) {
	filereader, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer filereader.Close()
	readerfunc := FindReader(filepath)
	game, err := readerfunc(filereader)
	if game != nil {
		game.Filename = filepath
	}
	return game, err
}

func FindReader(filepath string) func(io.Reader) (*Game, error) {
	if strings.HasSuffix(filepath, ".rle") || strings.HasSuffix(filepath, ".rle.txt") {
		return ReadRLE
	} else if strings.HasSuffix(filepath, ".life") || strings.HasSuffix(filepath, ".life.txt") {
		return ReadLife
	} else if strings.HasSuffix(filepath, ".cells") || strings.HasSuffix(filepath, ".cells.txt") {
		return ReadCells
	} else {
		return UnknownFiletypeReader
	}
}

func UnknownFiletypeReader(reader io.Reader) (*Game, error) {
	return nil, errors.New("Unsupported file type")
}

func ReadRLE(reader io.Reader) (*Game, error) {
	g := NewGame()

	bytes := make([]byte, 0, 1024)
	readBuf := make([]byte, 1024)
	for {
		n, err := reader.Read(readBuf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		if n > 0 {
			bytes = append(bytes, readBuf[:n]...)
		}
	}

	contents := string(bytes)
	lines := strings.Split(contents, "\n")

	var count_str strings.Builder
	count := func() int {
		if count_str.Len() == 0 {
			return 1
		} else {
			c, err := strconv.Atoi(count_str.String())
			check(err)
			count_str.Reset()
			return c
		}
	}

	var x, y, max_x Coord
	var expected_x, expected_y Coord
	rule := ""
	done := false

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "#N ") {
				g.Name = strings.TrimPrefix(line, "#N ")
			} else if strings.HasPrefix(line, "#O ") {
				g.Author = strings.TrimPrefix(line, "#O ")
			} else {
				g.Comments = append(g.Comments, strings.TrimPrefix(line, "#"))
			}
		} else {
			if strings.Contains(line, "=") {
				f := func(c rune) bool {
					return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '/'
				}
				fields := strings.FieldsFunc(line, f)
				if len(fields) < 6 || fields[0] != "x" || fields[2] != "y" || fields[4] != "rule" {
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
					rule = fields[5]
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
						return nil, errors.New(fmt.Sprintf("Got unknown code point %s", c))
					}
				}
			}
		}
		if done {
			break
		}
	}

	if strings.ToLower(rule) != "b3/s23" {
		return nil, errors.New("Unable to handle rule " + rule)
	}

	if !done {
		log.Println("WARN: Did not get terminator at end of RLE file")
	}

	if max_x < x {
		max_x = x
	}
	y += 1
	if max_x != expected_x || y != expected_y {
		log.Printf("WARN: Expected board of %dx%d got %dx%d", expected_x, expected_y, max_x, y)
	}

	return g, nil
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
	count  int
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
		case rel_x == last_x+1:
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

func (game *Game) SaveRLE(filepath string) error {
	fileWriter, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fileWriter.Close()
	return game.WriteRLE(fileWriter)
}

func (game *Game) WriteRLE(outfile io.Writer) error {

	min_cell, max_cell := game.Population.BoundingBox()
	if min_cell.X > max_cell.X {
		return nil
	}

	outwriter := bufio.NewWriter(outfile)
	if game.Name != "" {
		_, err := outwriter.WriteString(fmt.Sprintln("#N ", game.Name))
		if err != nil {
			return err
		}
	}
	if game.Author != "" {
		_, err := outwriter.WriteString(fmt.Sprintln("#O ", game.Author))
		if err != nil {
			return err
		}
	}
	for _, comment := range game.Comments {
		if !strings.HasPrefix(comment, "#") {
			_, err := outwriter.WriteString("#")
			if err != nil {
				return err
			}
		}
		_, err := outwriter.WriteString(comment)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(comment, "\n") {
			_, err := outwriter.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
	if len(game.Filename) > 0 {
		_, err := outwriter.WriteString(fmt.Sprintf("#C originally loaded from \"%s\"\n", game.Filename))
		if err != nil {
			return err
		}
	}
	if game.Generation > 0 {
		_, err := outwriter.WriteString(fmt.Sprintf("#C at generation %d\n", game.Generation))
		if err != nil {
			return err
		}
	}
	_, err := outwriter.WriteString(fmt.Sprintf("#C bounded by %d,%d -> %d,%d\n", min_cell.X, min_cell.Y, max_cell.X, max_cell.Y))
	if err != nil {
		return err
	}
	_, err = outwriter.WriteString(fmt.Sprintf("  x = %d, y = %d, rule = b3/s23\n", int(max_cell.X-min_cell.X+1), int(max_cell.Y-min_cell.Y+1)))
	if err != nil {
		return err
	}

	var line, newblob strings.Builder

	rle := game.ExtractRLE()
	for i := range rle {
		pair := rle[i]
		if pair.count > 1 {
			newblob.WriteString(fmt.Sprintf("%d%s", pair.count, pair.symbol))
		} else {
			newblob.WriteString(pair.symbol)
		}
		if line.Len()+newblob.Len() > max_line_length {
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
	if err != nil {
		return err
	}

	outwriter.Flush()

	return nil
}

func ReadLife(reader io.Reader) (*Game, error) {
	bytes := make([]byte, 0, 1024)
	readBuf := make([]byte, 1024)
	for {
		n, err := reader.Read(readBuf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		if n > 0 {
			bytes = append(bytes, readBuf[:n]...)
		}
	}

	content := string(bytes)
	lines := strings.Split(content, "\n")
	game := NewGame()
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
			case "\r":
			case " ":
			default:
				cells = append(cells, Cell{Coord(i), Coord(j)})
			}
		}
	}

	game.Population.Add(cells)

	return game, nil
}

func ReadCells(reader io.Reader) (*Game, error) {
	bytes := make([]byte, 0, 1024)
	readBuf := make([]byte, 1024)
	for {
		n, err := reader.Read(readBuf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		if n > 0 {
			bytes = append(bytes, readBuf[:n]...)
		}
	}

	content := string(bytes)
	lines := strings.Split(content, "\n")
	game := NewGame()
	cells := make(CellList, 0, 100)

	for lineNo := range lines {
		line := lines[lineNo]

		if strings.HasPrefix(line, "!") {
			game.Comments = append(game.Comments, line[1:])
			continue
		}

		chars := strings.Split(line, "")
		for x := range chars {
			c := chars[x]
			switch c {
			case "!":
				break
			case "O":
				cells = append(cells, Cell{Coord(x), Coord(lineNo)})
			case ".":
			case " ":
			case "\n":
			case "\r":
			default:
				return nil, errors.New(fmt.Sprintf("Unknown character '%s' in Cells file", c))
			}
		}
	}

	game.AddCells(cells)

	return game, nil
}
