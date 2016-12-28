package main

import (
	"fmt"
)

func main() {
	N := 5 // 3012210
	n := N

	// Give 1 present per elf.
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}

	i := 0
	for n > 1 {
		j := (i + n/2) % n
		fmt.Printf("Elves remaining: %v. Elf %d steals from elf %d. (i=%v, j=%v, n=%v)\n", p, p[i], p[j], i, j, n)
		if j == 0 {
			p = p[1:]
		} else if j == n-1 {
			p = p[0 : n-1]
		} else {
			p = append(p[0:j], p[j+1:]...)
		}

		n--
		if j > i {
			i++
		}
		i %= n
	}

	fmt.Printf("Elf %d ends up with all of the presents.\n", p[i])
}
