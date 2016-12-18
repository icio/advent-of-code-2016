package main

import (
	"fmt"
)

func main() {
	fav := uint(1350)
	orig := pos{1, 1}
	reachable := findOpen(fav, 50, orig)

	for y := uint(0); y < 23; y++ {
		for x := uint(0); x < 33; x++ {
			if step, reached := reachable[pos{x, y}]; reached {
				if step.step%10 == 0 {
					fmt.Print(step.step / 10)
				} else {
					fmt.Print(step.step % 10)
				}
			} else if wall(fav, x, y) {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Printf("Reachable: %d\n", len(reachable))
}

func findOpen(fav, steps uint, orig pos) map[pos]*step {
	origStep := &step{orig, 0}
	visited := map[pos]*step{orig: origStep}
	paths := []*step{origStep}

	for len(paths) > 0 {
		p := paths[0]
		paths = paths[1:]

		for dy := -1; dy < 2; dy += 1 {
			for dx := -1; dx < 2; dx += 1 {
				if (p.x == 0 && dx < 0) || (p.y == 0 && dy < 0) || (dx == 0) == (dy == 0) {
					continue
				}

				next := pos{uint(int(p.x) + dx), uint(int(p.y) + dy)}
				if _, seen := visited[next]; seen || wall(fav, next.x, next.y) {
					continue
				}

				nextStep := &step{next, p.step + 1}
				visited[next] = nextStep
				if nextStep.step < steps {
					paths = append(paths, nextStep)
				}
			}
		}
	}

	return visited
}

type step struct {
	pos
	step uint
}

type pos struct {
	x, y uint
}

func wall(fav, x, y uint) bool {
	w := x*x + 3*x + 2*x*y + y + y*y + fav
	n := 0
	for w > 0 {
		if w&1 == 1 {
			n++
		}
		w >>= 1
	}
	return n&1 == 1
}
