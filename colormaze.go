package main

import "fmt"
import "io/ioutil"
import "strings"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Point struct {
	row int
	col int
}

func (c Point) isNext(m *Maze) bool {
	// fmt.Printf("location: %v\n", c)
	// fmt.Printf("searching for: %c\n", m.sequence[m.seqIdx])
	// fmt.Printf("what I see: %c\n", byte(m.rows[c.row][c.col]))
	return m.sequence[m.seqIdx] == byte(m.rows[c.row][c.col])
}

// ensures the crumb won't be thrown over an edge
func (c Point) exists(width int, depth int) bool {
	return c.row >= 0 && c.row < depth && c.col >= 0 && c.col < width
}

// return adjacent points clockwise from last point
func (c Point) getAdjacent(direction Point, width int, depth int) []Point {
	possible := make([]Point, 4)

	var strafing Point
	if direction.row != 0 {
		strafing.row = 0
		strafing.col = -1
	} else {
		strafing.row = -1
		strafing.col = 0
	}

	// forward (continue)
	possible[0].row = c.row + direction.row
	possible[0].col = c.col + direction.col
	// left -- no ternaries??
	possible[1].row = c.row + strafing.row
	possible[1].col = c.col + strafing.col
	// right
	possible[2].row = c.row - strafing.row
	possible[2].col = c.col - strafing.col
	// back :(
	possible[3].row = c.row - direction.row
	possible[3].col = c.col - direction.col

	var actual []Point
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
	path     []Point // don't initialize this slice, just append to it
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
func (m *Maze) getEntryPoint() Point {
	entryRow := len(m.rows) - 1
	var crumb Point

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

// returns last crumb and direction we went
func (m *Maze) getLastPoint() (Point, Point) {
	last := m.path[len(m.path)-1]

	var direction Point
	if len(m.path) > 1 {
		direction.row = m.path[len(m.path)-2].row - last.row
		direction.row = m.path[len(m.path)-2].col - last.col
	} else {
		direction.row = 0
		direction.col = 0
	}

	return last, direction
}

// find next step (where sequence[n+1] == rows[len-1][row])
// probably want to prefer forward progress
func (m *Maze) getNextMove() Point {
	var c Point
	last, direction := m.getLastPoint()

	options := last.getAdjacent(direction, m.width, m.depth)

	for _, o := range options {
		if o.isNext(m) {
			c = o
			break
		}
	}

	return c
}

// adds the crumb to the path
// increments the index of the sequence so we know we good
func (m *Maze) dropPoint(c Point) {
	m.path = append(m.path, c)
	m.seqIdx++
	m.seqIdx %= len(m.sequence)
	fmt.Printf("dropped: %v\n", c)
}

func main() {
	maze := NewMaze("maze.txt")

	// find entry point (where sequence[0] == rows[len][row])
	current := maze.getEntryPoint()
	maze.dropPoint(current)

	// keep moving forward(ish) until we're at the top row
	iterations := 0
	for current.row != 0 && iterations < 20 {
		current := maze.getNextMove()
		maze.dropPoint(current)
		iterations++
	}

	fmt.Println("We're out of the woods!")
}

// sequence: O G
// 1: B O R O Y
// 2: O R B G R
// 3: B O G O Y
// 4: Y G B Y G
// 5: R O R B R
