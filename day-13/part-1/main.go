package main

import (
	"fmt"
	"log"
)

func main() {
	fav := uint(1350)
	path, err := findPath(fav, pos{1, 1}, pos{31, 39})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d steps: %#v\n", path.step, path)
}

func findPath(fav uint, orig, dest pos) (*step, error) {
	visited := map[pos]bool{orig: true}
	paths := []*step{&step{orig, 0, nil}}

	for len(paths) > 0 {
		p := paths[0]
		paths = paths[1:]

		for dy := -1; dy < 2; dy += 1 {
			for dx := -1; dx < 2; dx += 1 {
				if (p.x == 0 && dx < 0) || (p.y == 0 && dy < 0) || (dx == 0) == (dy == 0) {
					continue
				}

				next := pos{uint(int(p.x) + dx), uint(int(p.y) + dy)}
				if next == dest {
					return &step{next, p.step + 1, p}, nil
				}
				if wall(fav, next.x, next.y) {
					continue
				}
				if _, seen := visited[next]; seen {
					continue
				}
				log.Printf("Queueing %v -> %v", p.pos, next)
				paths = append(paths, &step{next, p.step + 1, p})
				visited[next] = true
			}
		}
	}

	return nil, fmt.Errorf("Could not find path from %v to %v", orig, dest)
}

type step struct {
	pos
	step uint
	prev *step
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
