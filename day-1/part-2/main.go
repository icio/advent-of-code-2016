package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	R = "R"
	L = "L"
)

type vec struct{ x, y int }

func main() {
	p := vec{0, 0}
	d := vec{0, 1}

	visited := make(map[vec]bool, 0)

	fmt.Fprintf(os.Stderr, "Starting at %#v.\n", p)

InstrLoop:
	for _, instr := range os.Args[1:] {
		visited[p] = true

		dir := strings.ToUpper(string(instr[0]))
		switch dir {
		case R:
			d = vec{d.y, -d.x}
		case L:
			d = vec{-d.y, d.x}
		default:
			panic(fmt.Sprintf("Unrecognised direction: %s", dir))
		}

		fwd, err := strconv.Atoi(strings.TrimRight(instr[1:], ","))
		if err != nil {
			panic(err)
		}

		for f := fwd; f > 0; f-- {
			p.x += d.x
			p.y += d.y

			if visited[p] {
				fmt.Fprintf(os.Stderr, "Already visited %v. HQ found.\n", p)
				break InstrLoop
			}
			visited[p] = true
		}

		fmt.Fprintf(os.Stderr, "Turned %s to face %#v and moved %d steps forward to %#v.\n", dir, d, fwd, p)
	}

	dist := abs(p.x) + abs(p.y)
	fmt.Println(dist)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}
