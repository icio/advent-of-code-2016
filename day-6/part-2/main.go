package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	var posCharMap []map[rune]charCount
	set := false

	for s.Scan() {
		if !set {
			set = true
			l := len(s.Text())

			posCharMap = make([]map[rune]charCount, l)
			for i := 0; i < l; i++ {
				posCharMap[i] = make(map[rune]charCount)
			}
		}
		for i, c := range s.Text() {
			t := posCharMap[i][c]
			t.char = c
			t.count++
			posCharMap[i][c] = t // TODO: How can we define t s.t. this isn't required?
		}
	}

	msg := ""
	for _, charMap := range posCharMap {
		counts := make(charList, 0)
		for _, count := range charMap {
			counts = append(counts, count)
		}

		sort.Sort(counts)
		fmt.Println(counts)
		msg += string(counts[0].char)
	}

	fmt.Println(msg)
}

type charCount struct {
	char  rune
	count int
}

type charList []charCount

func (c charList) Len() int      { return len(c) }
func (c charList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c charList) Less(i, j int) bool {
	return c[i].count < c[j].count
}
