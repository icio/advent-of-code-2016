package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	SAFE   = '.'
	UNSAFE = '^'
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	n := 0
	if len(os.Args) > 1 {
		n, _ = strconv.Atoi(os.Args[1])
	}
	if n == 0 {
		n = 40
	}
	fmt.Println("Safe tiles:", safeTiles(string(input), n))
}

func safeTiles(input string, rows int) (safe int) {
	above := []rune(input)
	for _, c := range above {
		if c == SAFE {
			safe++
		}
	}
	fmt.Fprintln(os.Stderr, input)

	n := len(input)
	for row := 1; row < rows; row++ {
		nextRow := make([]rune, n)
		for x := 0; x < n; x++ {
			if x == 0 {
				if above[x+1] == UNSAFE {
					nextRow[x] = UNSAFE
				} else {
					nextRow[x] = SAFE
					safe++
				}
				continue
			}
			if x == n-1 {
				if above[x-1] == UNSAFE {
					nextRow[x] = UNSAFE
				} else {
					nextRow[x] = SAFE
					safe++
				}
				continue
			}
			if above[x-1] != above[x+1] && (above[x-1] == UNSAFE || above[x+1] == UNSAFE) {
				nextRow[x] = UNSAFE
			} else {
				nextRow[x] = SAFE
				safe++
			}
		}
		fmt.Fprintln(os.Stderr, string(nextRow))
		above = nextRow
	}
	return
}
