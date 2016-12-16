package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// From the problem description in layout.txt.
	f := Factory(
		0<<Elevator |
			(Strontium<<Generator | Strontium | Plutonium<<Generator | Plutonium) |
			(Thulium<<Generator|Ruthenium<<Generator|Ruthenium|Curium<<Generator|Curium)<<(1*FloorSize) |
			(Thulium)<<(2*FloorSize),
	)

	if false {
		// From the problem description in example.txt.
		f = Factory(
			0<<Elevator |
				(Strontium|Plutonium)<<(0*FloorSize) |
				(Strontium<<Generator)<<(1*FloorSize) |
				(Plutonium<<Generator)<<(2*FloorSize),
		)
	}

	completePath, complete := CompletePath(f)
	if !complete {
		panic("Failed to find path to complete state.")
	}

	for completePath != (*path)(nil) {
		fmt.Println(completePath.Steps)
		fmt.Println(completePath.State.String())
		completePath.State.Display(os.Stdout)
		completePath = completePath.Previous
	}
}

// CompletePath performs a breath-first search through safe movements of
// factory state from the given to the first complete state.
func CompletePath(f Factory) (*path, bool) {
	visited := map[Factory]bool{
		f: true,
	}
	reachable := []*path{
		&path{f, nil, 0},
	}

	for len(reachable) > 0 {
		// Unshift an item from the path queue.
		p := reachable[0]
		reachable = reachable[1:]

		for _, next := range p.State.Moves(true) {
			if _, seen := visited[next]; seen {
				continue
			}
			nextP := &path{
				State:    next,
				Previous: p,
				Steps:    p.Steps + 1,
			}
			if next.Complete() {
				return nextP, true
			}
			visited[next] = true
			reachable = append(reachable, nextP)
		}
	}

	return nil, false
}

type path struct {
	State    Factory
	Previous *path
	Steps    uint
}

func debugMasks() {
	fmt.Printf("%064b\n", FullMask)

	for floor := uint64(0); floor < Floors; floor++ {
		floorOffset := floor * FloorSize
		fmt.Printf("%064b Floor #%d\n", FloorMask<<floorOffset, floor)
		fmt.Printf("%064b - chips\n", ChipMask<<floorOffset)
		fmt.Printf("%064b - generators\n", GeneratorMask<<floorOffset)
	}
}

const (
	Floors     uint64 = 4
	Components        = 5

	ComponentNames = "SPTRC"
	TypeNames      = "CG"

	Generator uint64 = Components
	Elevator         = Floors * FloorSize
	FloorSize        = Components * 2
)

const (
	Strontium uint64 = 1 << iota
	Plutonium
	Thulium
	Ruthenium
	Curium
)

const (
	FullMask      = ^uint64(0)
	FloorsMask    = FullMask >> (64 - Floors*FloorSize)
	ElevatorMask  = FullMask ^ FloorsMask
	ChipMask      = FullMask >> (64 - Components)
	GeneratorMask = ChipMask << Components
	FloorMask     = ChipMask | GeneratorMask
)

type Factory uint64

func (f Factory) String() string {
	return fmt.Sprintf("%064b", f)
}

// Elevator identifies the floor on which you can find the lift.
func (f Factory) Elevator() uint64 {
	return uint64(f >> Elevator)
}

// Safe returns whether all chips are adequately shielded from any generators on the same floor.
func (f Factory) Safe() bool {
	m := ChipMask
	g := uint64(f) >> Components

	for l := uint64(0); l < Floors; l++ {
		if g&m != 0 && uint64(f)&m&^g != 0 {
			return false
		}
		m <<= FloorSize
	}

	return true
}

// Complete returns whether all components in the Factory are on the top floor.
func (f Factory) Complete() bool {
	g := uint64(f) & FloorsMask
	top := g & (FloorMask << (FloorSize * (Floors - 1)))
	return g == top
}

// Moves returns the factory states which can be reached by transporting 1 or 2
// components from the current floor to a neighbouring floor.
func (f Factory) Moves(onlySafe bool) (moves []Factory) {
	E := f.Elevator()
	var startE uint64
	if E > 0 {
		startE = E - 1
	} else {
		startE = E + 1
	}

	g := uint64(f)
	l := g >> (E * FloorSize)
	for a := uint64(1); a < FloorMask; a <<= 1 {
		if l&a == 0 {
			continue
		}
		for b := a; b < FloorMask; b <<= 1 {
			if l&b == 0 {
				continue
			}

			// The components to be moved from floor E to e.
			c := l & (a | b)

			for e := startE; e <= E+1 && e < Floors; e += 2 {
				move := Factory(e<<Elevator | (g&FloorsMask)&^(c<<(E*FloorSize)) | c<<(e*FloorSize))
				if onlySafe && !move.Safe() {
					continue
				}
				moves = append(moves, move)
			}
		}
	}

	return
}

// Display prints the factory layout into w.
func (f Factory) Display(w io.Writer) {
	e := f.Elevator()

	for l := Floors - 1; l < Floors; l-- {
		g := uint64(f) >> (FloorSize * l)

		fmt.Printf("%02d ", l)
		if l == e {
			fmt.Print("E | ")
		} else {
			fmt.Print("  | ")
		}

		for t := 1; t >= 0; t-- {
			for c := Components - 1; c >= 0; c-- {
				if g&(1<<uint(c+t*Components)) > 0 {
					fmt.Fprint(w, string(TypeNames[t]), string(ComponentNames[c]), " ")
				} else {
					fmt.Fprint(w, ".  ")
				}
			}
			fmt.Fprint(w, " ")
		}
		fmt.Fprintln(w)
	}
}
