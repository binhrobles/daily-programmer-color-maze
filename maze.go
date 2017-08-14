package main

import "io/ioutil"
import "strings"
import "fmt"

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
		direction.row = last.row - m.path[len(m.path)-2].row
		direction.col = last.col - m.path[len(m.path)-2].col
	} else {
		// point em north
		direction.row = -1
		direction.col = 0
	}

	return last, direction
}

// find next step (where sequence[n+1] == rows[len-1][row])
// probably want to prefer forward progress
func (m *Maze) getNextMove() Point {
	var c Point
	c.row = -1

	last, direction := m.getLastPoint()

	options := last.getAdjacent(direction, m.width, m.depth)

	for _, o := range options {
		if o.isNext(m) {
			c = o
			// TODO: queue this option up, somehow
			// need to create a queue of choices that we can refer back to
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
