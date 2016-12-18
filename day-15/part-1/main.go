package main

import (
	"fmt"
)

func main() {
	// The example:
	// disks := []Disk{
	// 	{1, 5, 4},
	// 	{2, 2, 1},
	// }

	// Copy-pasta from input.txt
	disks := []Disk{
		{1, 17, 15},
		{2, 3, 2},
		{3, 19, 4},
		{4, 13, 2},
		{5, 7, 2},
		{6, 5, 0},
		{7, 11, 0}, // Part 2
	}
	fmt.Println(firstOpenTime(disks))
}

func firstOpenTime(disks []Disk) uint {
Time:
	for t := uint(0); ; t++ {
		for n, d := range disks {
			if !d.Open(t + uint(n)) {
				continue Time
			}
		}

		return t
	}
}

type Disk struct {
	Index  uint
	Holes  uint
	Offset uint
}

func (d Disk) Open(time uint) bool {
	return (time+d.Offset+1)%d.Holes == 0
}
