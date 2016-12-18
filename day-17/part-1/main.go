package main

import (
	"crypto/md5"
	"fmt"
	"os"
)

func main() {
	fmt.Println(findPath(os.Args[1]).moves)
}

var dirs [4]dir

func init() {
	dirs = [4]dir{
		{0, -1, "U"},
		{0, 1, "D"},
		{-1, 0, "L"},
		{1, 0, "R"},
	}
}

func findPath(seed string) path {
	paths := []path{path{"", 0, 0}}

	for len(paths) > 0 {
		p := paths[0]
		paths = paths[1:]

		doors := fmt.Sprintf("%x", md5.Sum([]byte(seed+p.moves)))
		for c := 0; c < 4; c++ {
			if doors[c] < 'b' {
				continue
			}

			d := dirs[c]
			t := path{p.moves + d.m, p.x + d.x, p.y + d.y}
			if t.x == 3 && t.y == 3 {
				return t
			}
			if t.x >= 0 && t.x <= 3 && t.y >= 0 && t.y <= 3 {
				paths = append(paths, t)
			}
		}
	}

	panic("Could not find a path!")
}

func success(p path) bool {
	return p.x == 3 && p.y == 3
}

type dir struct {
	x, y int
	m    string
}

type path struct {
	moves string
	x, y  int
}
