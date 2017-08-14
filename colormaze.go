package main

import "fmt"
import "flag"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	filename := flag.String("maze", "maze.txt", "the file containing the maze sequence and definition")

	flag.Parse()
	maze := NewMaze(*filename)

	// find entry point (where sequence[0] == rows[len][row])
	current := maze.getEntryPoint()
	maze.dropPoint(current)

	// keep moving forward(ish) until we're at the top row
	iterations := 0
	for current.row != 0 && iterations < maze.width*maze.depth {
		current = maze.getNextMove()

		// we weren't able to go down this path
		if current.row == -1 {
			// go back to last point that had another valid option
		} else {
			maze.dropPoint(current)
		}
		iterations++
	}

	fmt.Println("We're out of the woods!")
}
