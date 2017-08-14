package main

import "fmt"
import "io/ioutil"
import "strings"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// making the assumption that the maze is a square
func getMazeSpecs(s string) (string, []string) {
	f, err := ioutil.ReadFile(s)
	check(err)

	raw := strings.Split(string(f), "\n")
	replacer := strings.NewReplacer(" ", "")

	sequence := replacer.Replace(raw[0])
	width := len(replacer.Replace(raw[1]))

	fmt.Printf("width: %d depth: %d\n", width, len(raw)-2)

	rows := make([]string, len(raw)-2)

	for i := 0; i < len(raw)-2; i++ {
		rows[i] = replacer.Replace(raw[i+1])
	}

	return sequence, rows
}

func main() {
	sequence, rows := getMazeSpecs("maze.txt")
	fmt.Printf("sequence: %s\nmaze: \n%v", sequence, rows)

	// TODO: find entry point (where sequence[0] == rows[len][x])

	// TODO: find next step (where sequence[n+1] == rows[len-1][x])
	// probably want to prefer forward progress
}

// sequence: O G
// 1: B O R O Y
// 2: O R B G R
// 3: B O G O Y
// 4: Y G B Y G
// 5: R O R B R
