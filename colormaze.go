package main

import "fmt"
import "io/ioutil"
import "strings"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Crumb struct {
	row int
	col int
}

func (c Crumb) isNext(m *Maze) bool {
	return m.sequence[m.seqIdx] == byte(m.rows[c.row][c.col])
}

// ensures the crumb won't be thrown over an edge
func (c Crumb) exists(width int, depth int) bool {
	return c.row >= 0 && c.row < depth && c.col >= 0 && c.col < width
}

// return adjacent points north, west, east, then south
func (c Crumb) getAdjacent(width int, depth int) []Crumb {
	possible := make([]Crumb, 4)

	// north
	possible[0].row = c.row - 1
	possible[0].col = c.col
	// west
	possible[1].row = c.row
	possible[1].col = c.col - 1
	// east
	possible[2].row = c.row
	possible[2].col = c.col + 1
	// south :(
	possible[3].row = c.row + 1
	possible[3].col = c.col

	var actual []Crumb
	for _, p := range possible {
		if p.exists(width, depth) {
			actual = append(actual, p)
		}
	}

	return actual
}

type Maze struct {
	sequence string
	seqIdx   int
	rows     []string
	path     []Crumb // don't initialize this slice, just append to it
	width    int
	depth    int
}

// making the assumption that the maze is a square
func NewMaze(s string) *Maze {
	f, err := ioutil.ReadFile(s)
	check(err)

	raw := strings.Split(string(f), "\n")
	replacer := strings.NewReplacer(" ", "")

	sequence := replacer.Replace(raw[0])
	width := len(replacer.Replace(raw[1]))
	depth := len(raw) - 2

	fmt.Println("Maze Specs\n----------")
	fmt.Printf("width: %d depth: %d\n", width, depth)
	fmt.Printf("sequence: %s\n", sequence)

	rows := make([]string, len(raw)-2)

	for i := 0; i < len(raw)-2; i++ {
		rows[i] = replacer.Replace(raw[i+1])
		fmt.Printf("%s\n", raw[i+1])
	}

	var m Maze
	m.sequence = sequence
	m.seqIdx = 0
	m.rows = rows
	m.width = width
	m.depth = depth
	return &m
}

// find entry point (where sequence[0] == rows[len][row])
// the bottom row (highest index) is the beginning
func (m *Maze) getEntryPoint() Crumb {
	entryRow := len(m.rows) - 1
	var crumb Crumb

	// iterate through every char in the entry row
	// until we find the char we want
	for i, c := range m.rows[entryRow] {
		if m.sequence[m.seqIdx] == byte(c) {
			crumb.row = entryRow
			crumb.col = i
			break
		}
	}

	return crumb
}

func (m *Maze) getLastCrumb() Crumb {
	return m.path[len(m.path)-1]
}

func (m *Maze) goForward() Crumb {
	var c Crumb
	last := m.getLastCrumb()

	options := last.getAdjacent(m.width, m.depth)

	for _, o := range options {
		if o.isNext(m) {
			c = o
		}
	}

	fmt.Printf("c: %v\n", c)

	return c
}

// adds the crumb to the path
// increments the index of the sequence so we know we good
func (m *Maze) dropCrumb(c Crumb) {
	m.path = append(m.path, c)
	m.seqIdx++
	m.seqIdx %= len(m.sequence)
	fmt.Printf("Path forward: %v\n", m.path)
}

func main() {
	maze := NewMaze("maze.txt")

	// find entry point (where sequence[0] == rows[len][row])
	entry := maze.getEntryPoint()
	maze.dropCrumb(entry)

	// TODO: find next step (where sequence[n+1] == rows[len-1][row])
	// probably want to prefer forward progress
	maze.goForward()
}

// sequence: O G
// 1: B O R O Y
// 2: O R B G R
// 3: B O G O Y
// 4: Y G B Y G
// 5: R O R B R
