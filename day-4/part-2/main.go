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

	for s.Scan() {
		room := s.Text()
		valid, name, sectorId := roomValid(room)

		if !valid {
			continue
		}

		fmt.Println(decrypt(name, sectorId), sectorId)
	}

	if err := s.Err(); err != nil {
		panic(err)
	}
}

func roomValid(room string) (bool, string, int) {
	name := room[0 : len(room)-11]
	sectorId, _ := strconv.Atoi(room[len(room)-10 : len(room)-7])
	checkSum := room[len(room)-6 : len(room)-1]

	return validateChecksum(name, checkSum), name, sectorId
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

func decrypt(name string, sectorId int) string {
	d := ""
	for _, c := range name {
		if c == '-' {
			d += " "
			continue
		}
		d += string(97 + ((int(c)-97)+sectorId+26*4)%26)
	}

	return d
}
