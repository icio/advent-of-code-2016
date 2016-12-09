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

	p := vec{1, 1}
	keypad := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	for _, instr := range os.Args[1:] {
		for _, d := range instr {
			dir, ok := dirs[d]
			if !ok {
				panic(fmt.Sprintf("Could not parse direction: %v", d))
			}

			p.x = lim(p.x + dir.x)
			p.y = lim(p.y + dir.y)
		}
		fmt.Print(keypad[p.y][p.x])
	}

	fmt.Println()
}

func lim(p int) int {
	if p < 0 {
		return 0
	}
	if p > 2 {
		return 2
	}
	return p
}
