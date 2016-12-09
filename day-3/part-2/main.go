package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

Read:
	for {
		var tris [3][3]string
		for i := 0; i < 3; i++ {
			for t := 0; t < 3; t++ {
				if !scanner.Scan() {
					if t == 0 && i == 0 {
						break Read
					}
					panic("EOF reached unexpectedly.")
				}
				tris[t][i] = scanner.Text()
			}
		}

		for _, tri := range tris {
			valid, err := validTriangle(tri[0], tri[1], tri[2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			} else if valid {
				fmt.Println(tri)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func validTriangle(strA, strB, strC string) (bool, error) {
	// Parse the three numbers out
	a, err_a := strconv.Atoi(strA)
	b, err_b := strconv.Atoi(strB)
	c, err_c := strconv.Atoi(strC)

	if err_a != nil || err_b != nil || err_c != nil {
		return false, fmt.Errorf("%s %s %s", err_a, err_b, err_c)
	}

	// Find the longest side.
	var hyp int
	if a < b {
		if b < c {
			hyp = c
		} else {
			hyp = b
		}
	} else {
		if a < c {
			hyp = c
		} else {
			hyp = a
		}
	}

	return a+b+c-hyp > hyp, nil
}
