package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) < 15 {
			panic(fmt.Sprintf("Expected 3 numbers right-aligned at columns 5, 10, and 15. %q invalid.", t))
		}

		// Parse the three numbers out
		a, err_a := strconv.Atoi(strings.Trim(t[0:5], " "))
		b, err_b := strconv.Atoi(strings.Trim(t[5:10], " "))
		c, err_c := strconv.Atoi(strings.Trim(t[10:15], " "))

		if err_a != nil || err_b != nil || err_c != nil {
			panic(fmt.Errorf("%s %s %s", err_a, err_b, err_c))
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

		if a+b+c-hyp > hyp {
			fmt.Println(t)
		}
	}
	if scanner.Err() != nil {
		panic(scanner.Err)
	}
}
