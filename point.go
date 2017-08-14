package main

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
