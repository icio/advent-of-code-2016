package main

import (
	"fmt"
)

type elf struct {
	name int
	out  bool
	next int
}

func main() {
	N := 3012210
	n := N

	// Give 1 present per elf.
	p := make([]elf, n)
	for i := 0; i < n; i++ {
		e := &p[i]
		e.name = i + 1
		e.next = (i + 1) % n
	}

	i := 0
	for n > 1 {
		curr := &p[i]

		// Find the next elf to take out of the game.
		next := &p[curr.next]
		for next.out {
			next = &p[next.next]
		}

		next.out = true
		n--

		// Point the current elf to the next remaining elf.
		for next.out {
			next = &p[next.next]
		}

		// Progress to the next elf to take a shot.
		i = next.name - 1
		curr.next = i
	}

	fmt.Printf("i=%v, n=%v, p[i]=%+v\n", i, n, p[i])
	fmt.Printf("Elf %d ends up with all of the presents.\n", p[i].name)
}
