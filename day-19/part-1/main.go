package main

import (
	"fmt"
)

func main() {
	n := 3012210

	// Give 1 present per elf.
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = 1
	}

	i := 0
	for {
		if p[i] == 0 {
			i = (i + 1) % n
			continue
		}

		// Find elf j to steal from.
		j := (i + 1) % n
		for ; j != i; j = (j + 1) % n {
			if p[j] > 0 {
				break
			}
		}

		// Check if we went full circle.
		if j == i {
			panic(fmt.Errorf("%d is the only elf left. Should have returned elsewhere.", i+1))
		}

		// Steal presents for i from j.
		fmt.Printf("Elf %d steals %d presents from Elf %d.\n", i+1, p[j], j+1)
		p[i], p[j] = p[i]+p[j], 0

		if p[i] == n {
			fmt.Printf("All %d presents end up with Elf %d.\n", p[i], i+1)
			return
		}

		i = (j + 1) % n
	}
}
