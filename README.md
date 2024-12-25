# John Conway's Game of Life logic 
## implemented in Go
## by Mitch Patenaude <mitch@mitchpatenaude.net>
## Copyright 2024

This is a simple, but fairly performant and memory efficient implementation 
of [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life)

This is simply a library with no user interface.  The main data structure is the
**Game** structure.  It has a member element of type **Population**, which is a 
collection of **Cell**s, which are themselves just a pair of X and Y **Coord**s.

For now, a **Coord** is just a int64 type, and that seems the most flexible, but 
for memory efficiency could be changed to an int32.

```
type Coord int64

type Cell struct {
    X, Y Coord
}

type Population map[Cell]bool

type Game struct {
    Filename string
    Population Population
    History []Population
    HistorySize int
    Comments []string
    Generation int
}

```
The Only interesting methods on the **Population** type are
```
func (pop Population) Add(new_cells []Cell) 
```
Adds new cells to the population

and
```
func (current Population) Step() Population 
```
Calculates the next generation of the population.

```
func NewGame() *Game 
```

```
func Load(filepath string) *Game
```

```
func (game *Game) SetHistorySize(size int)
```

```
func (game *Game) AddCell(cell Cell)
```

```
func (game *Game) AddCells(cells []Cell)
```

```
func (game *Game) Next()
```


