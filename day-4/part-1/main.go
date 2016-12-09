package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	t := 0

	for s.Scan() {
		room := s.Text()
		valid, sectorId := roomValid(room)
		if !valid {
			continue
		}

		fmt.Fprintln(os.Stderr, room)
		t += sectorId
	}

	if err := s.Err(); err != nil {
		panic(err)
	}

	fmt.Println(t)
}

func roomValid(room string) (bool, int) {
	name := room[0 : len(room)-11]
	sectorId, _ := strconv.Atoi(room[len(room)-10 : len(room)-7])
	checkSum := room[len(room)-6 : len(room)-1]

	return validateChecksum(name, checkSum), sectorId
}

type charCount struct {
	char  rune
	count int
}

type charList []charCount

func (c charList) Len() int      { return len(c) }
func (c charList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (c charList) Less(i, j int) bool {
	if c[i].count == c[j].count && c[i].char < c[j].char {
		return true
	}
	return c[i].count > c[j].count
}

func validateChecksum(name, checkSum string) bool {
	charMap := make(map[rune]charCount)
	for _, c := range name {
		if c == '-' {
			continue
		}

		t := charMap[c]
		t.char = c
		t.count++
		charMap[c] = t
	}

	counts := make(charList, 0, len(charMap))
	for _, count := range charMap {
		counts = append(counts, count)
	}

	sort.Sort(counts)

	cs := ""
	for i := 0; i < 5; i++ {
		cs = cs + string(counts[i].char)
	}

	return checkSum == cs
}
