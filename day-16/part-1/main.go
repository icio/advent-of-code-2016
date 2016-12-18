package main

import (
	"fmt"
	"os"
)

func main() {
	contents := "00101000101111010"
	N := 272
	if len(os.Args) > 1 {
		// Part 2
		N = 35651584
	}

	for len(contents) < N {
		contents = grow(contents)
	}
	contents = contents[:N]

	if N < 1000 {
		fmt.Println(contents)
	}
	fmt.Println("Checksum:", checksum(contents))
}

func grow(input string) string {
	n := len(input)
	flip := []rune(input)
	if n%2 == 1 {
		flip[n/2] = inv(flip[n/2])
	}
	for i := 0; i < n/2; i++ {
		flip[i], flip[n-1-i] = inv(flip[n-1-i]), inv(flip[i])
	}

	return input + "0" + string(flip)
}

func inv(c rune) rune {
	if c == '0' {
		return '1'
	}
	return '0'
}

func checksum(input string) string {
	for len(input)&1 == 0 {
		chk := make([]rune, len(input)/2)
		for c := 0; c < len(input); c += 2 {
			if input[c] == input[c+1] {
				chk[c/2] = '1'
			} else {
				chk[c/2] = '0'
			}
		}
		input = string(chk)
	}
	return input
}
