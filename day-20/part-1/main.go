package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	ips := NewBlockedIPRange()

	// Read the blocked ranges.
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		from, to, err := parseRange(s.Text())
		if err != nil {
			panic(err)
		}
		ips.Block(from, to)
	}
	if err := s.Err(); err != nil {
		panic(err)
	}

	fmt.Println(ips.LowestAllowed())
}

func parseRange(r string) (int, int, error) {
	fromTo := strings.SplitN(r, "-", 2)
	if len(fromTo) != 2 {
		return 0, 0, fmt.Errorf("expected range of format '%d-%d' but got %q", r)
	}

	from, err := strconv.Atoi(fromTo[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid range from in %q: %s", r, err)
	}

	to, err := strconv.Atoi(fromTo[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid range to in %q: %s", r, err)
	}

	return from, to, nil
}

type BlockedIPRange struct {
	ranges [][2]int
}

func NewBlockedIPRange() *BlockedIPRange {
	return &BlockedIPRange{}
}

func (b *BlockedIPRange) Block(from, to int) {
	fmt.Printf("Blocking from %d to %d.\n", from, to)
	b.ranges = append(b.ranges, [2]int{from, to})
}

func (b *BlockedIPRange) LowestAllowed() int {
	allowed := 0
Search:
	for {
		for _, bounds := range b.ranges {
			if allowed >= bounds[0] && allowed <= bounds[1] {
				allowed = bounds[1] + 1
				continue Search
			}
		}
		break
	}
	return allowed
}

func (b *BlockedIPRange) Allowed(n int) bool {
	return true
}
