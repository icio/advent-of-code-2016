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

func main() {
	x, y := 0, 0
	dx, dy := 0, 1

	for _, instr := range os.Args[1:] {
		dir := strings.ToUpper(string(instr[0]))
		switch dir {
		case R:
			dx, dy = dy, -dx
		case L:
			dx, dy = -dy, dx
		default:
			panic(fmt.Sprintf("Unrecognised direction: %s", dir))
		}

		fwd, err := strconv.Atoi(strings.TrimRight(instr[1:], ","))
		if err != nil {
			panic(err)
		}

		x += fwd * dx
		y += fwd * dy

		fmt.Fprintf(os.Stderr, "Turned %s to face (dx: %d, dy: %d) and moved %d steps forward to (x: %d, y: %d)\n", dir, dx, dy, fwd, x, y)
	}

	dist := abs(x) + abs(y)
	fmt.Println(dist)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}
