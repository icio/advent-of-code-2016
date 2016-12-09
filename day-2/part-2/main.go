package main

import (
	"fmt"
	"os"
)

type vec struct{ x, y int }

func main() {
	dirs := map[rune]vec{
		'U': {0, -1},
		'D': {0, 1},
		'L': {-1, 0},
		'R': {1, 0},
	}

	p := vec{1, 3}
	keypad := []string{
		"       ",
		"   1   ",
		"  234  ",
		" 56789 ",
		"  ABC  ",
		"   D   ",
		"       ",
	}

	for _, instr := range os.Args[1:] {
		for _, d := range instr {
			dir, ok := dirs[d]
			if !ok {
				panic(fmt.Sprintf("Could not parse direction: %v", d))
			}

			if keypad[p.y+dir.y][p.x+dir.x] == ' ' {
				continue
			}

			p.x += dir.x
			p.y += dir.y
		}
		fmt.Print(string(keypad[p.y][p.x]))
	}

	fmt.Println()
}
