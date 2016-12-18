package main

import (
	"crypto/md5"
	"fmt"
	"os"
)

func main() {
	longest := findLongestPath(os.Args[1])
	fmt.Println(longest.moves)
	fmt.Println(len(longest.moves))
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

func findLongestPath(seed string) path {
	longest := path{"", 0, 0}
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
				if len(t.moves) > len(longest.moves) {
					longest = t
				}
				continue
			}
			if t.x >= 0 && t.x <= 3 && t.y >= 0 && t.y <= 3 {
				paths = append(paths, t)
			}
		}
	}

	return longest
}

type dir struct {
	x, y int
	m    string
}

type path struct {
	moves string
	x, y  int
}
