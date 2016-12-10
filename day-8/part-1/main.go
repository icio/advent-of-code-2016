package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	d := NewDisplay()

	for s.Scan() {
		fmt.Fprintln(os.Stderr, s.Text())

		instruction, a, b, err := parseInstruction(s.Text())
		if err != nil {
			panic(err)
		}

		switch instruction {
		case INSTR_FILL:
			d.Fill(a, b)
		case INSTR_ROT_ROW:
			d.RotateRow(a, b)
		case INSTR_ROT_COL:
			d.RotateColumn(a, b)
		}

		d.Print(os.Stderr)
	}

	fmt.Println("Lit: ", d.Count())
	d.Print(os.Stdout)

	if err := s.Err(); err != nil {
		panic(err)
	}
}

const (
	INSTR_FILL int = iota
	INSTR_ROT_ROW
	INSTR_ROT_COL
)

func parseInstruction(t string) (inst, a, b int, err error) {
	if strings.HasPrefix(t, "rect ") {
		inst = INSTR_FILL
		a = parseInt(t[5:])
		b = parseInt(t[strings.Index(t, "x")+1:])
	} else if strings.HasPrefix(t, "rotate row y=") {
		inst = INSTR_ROT_ROW
		a = parseInt(t[13:])
		b = parseInt(t[strings.Index(t, "by ")+3:])
	} else if strings.HasPrefix(t, "rotate column x=") {
		inst = INSTR_ROT_COL
		a = parseInt(t[16:])
		b = parseInt(t[strings.Index(t, "by ")+3:])
	} else {
		err = fmt.Errorf("Unrecognised instruction: %q", t)
	}

	return
}

func parseInt(t string) int {
	digits := 0
	for i := range t {
		if t[i] >= 48 && t[i] <= 57 {
			digits++
		} else {
			break
		}
	}

	if digits == 0 {
		return 0
	}
	v, _ := strconv.Atoi(t[:digits])
	return v
}

const (
	HEIGHT int = 6
	WIDTH  int = 50
)

type Display [HEIGHT][WIDTH]int

func NewDisplay() *Display {
	var d Display
	return &d
}

func (d *Display) Print(w io.Writer) {
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if d[y][x]&1 == 1 {
				w.Write([]byte("#"))
			} else {
				w.Write([]byte("."))
			}
		}
		w.Write([]byte("\n"))
	}
	w.Write([]byte("\n"))
}

func (d *Display) RotateColumn(x, dy int) {
	for y := 0; y < HEIGHT; y++ {
		py := (y - dy + 3*HEIGHT) % HEIGHT
		p := d[py][x]
		if py < y {
			p >>= 1
		}
		d[y][x] = (d[y][x] << 1) + (p & 1)
	}
}

func (d *Display) RotateRow(y, dx int) {
	for x := 0; x < WIDTH; x++ {
		px := (x - dx + 3*WIDTH) % WIDTH
		p := d[y][px]
		if px < x {
			p >>= 1
		}
		d[y][x] = (d[y][x] << 1) + (p & 1)
	}
}

func (d *Display) Fill(w, h int) {
	for y := 0; y < h && y < HEIGHT; y++ {
		for x := 0; x < w && x < WIDTH; x++ {
			d[y][x] = 1
		}
	}
}

func (d *Display) Count() (count int) {
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if d[y][x]&1 == 1 {
				count++
			}
		}
	}
	return
}
