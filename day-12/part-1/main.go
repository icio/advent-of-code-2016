package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const debug = false

func main() {
	program, err := parse(os.Stdin)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:")
	memory := make(MapMemory)
	run(memory, program)
	dump(memory)

	fmt.Println("Part 2:")
	memory = make(MapMemory)
	memory.Write("c", 1)
	run(memory, program)
	dump(memory)
}

func run(memory Memory, program []Instruction) {
	cursor := 0
	for cursor >= 0 && cursor < len(program) {
		inst := program[cursor]
		if debug {
			log.Printf("%d: %s", cursor, inst)
		}
		cursor += 1 + inst.Apply(memory)
		if debug {
			log.Println(memory.String())
		}
	}
}

func dump(memory MapMemory) {
	for register, value := range memory {
		fmt.Printf("%s\t%d\n", register, value)
	}
}

func isInt(s string) bool {
	for i, c := range s {
		if i == 0 && c == '-' {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func parse(r io.Reader) ([]Instruction, error) {
	inst := make([]Instruction, 0)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := strings.Split(s.Text(), " ")
		if len(t) < 2 {
			return nil, fmt.Errorf("Input too short: %q", s.Text())
		}
		switch t[0] {
		case "cpy":
			if len(t) != 3 || t[1] == "" || t[2] == "" {
				return nil, fmt.Errorf("expected `cpy (REG|VAL) REG` but got %q", s.Text())
			}
			var reg string
			var val int
			if isInt(t[1]) {
				val, _ = strconv.Atoi(t[1])
			} else {
				reg = t[1]
			}
			inst = append(inst, &Copy{
				fromRegister: reg,
				fromValue:    val,
				toRegister:   t[2],
			})
		case "inc":
			if len(t) != 2 || t[1] == "" {
				return nil, fmt.Errorf("expected `inc REG` but got %q", s.Text())
			}
			inst = append(inst, &Add{
				register: t[1],
				d:        1,
			})
		case "dec":
			if len(t) != 2 || t[1] == "" {
				return nil, fmt.Errorf("expected `dec REG` but got %q", s.Text())
			}
			inst = append(inst, &Add{
				register: t[1],
				d:        -1,
			})
		case "jnz":
			reg := t[1]
			if isInt(reg) {
				reg = ""
			}
			if !isInt(t[2]) {
				return nil, fmt.Errorf("expected `jnz (REG|VAL) INT` but got %q", s.Text())
			}
			n, _ := strconv.Atoi(t[2])
			inst = append(inst, &JumpNonZero{
				register: reg,
				offset:   n,
			})
		default:
			panic(fmt.Sprintf("Unrecognised instruction: %q", t))
		}
	}

	return inst, nil
}

type Memory interface {
	Read(register string) int
	Write(register string, value int)
	String() string
}

type MapMemory map[string]int

func (m MapMemory) Read(register string) int {
	return m[register]
}

func (m MapMemory) Write(register string, value int) {
	m[register] = value
}

func (m MapMemory) String() string {
	return fmt.Sprintf("%+v", (map[string]int)(m))
}

type Instruction interface {
	// Apply enacts the operation contained in the instruction, mutating the
	// memory and returning the number of instructions by which the cursor
	// should jump (0 to progress as normal; -1 to repeat forever).
	Apply(Memory) int
	String() string
}

type Copy struct {
	fromValue    int
	fromRegister string
	toRegister   string
}

func (c *Copy) Apply(m Memory) (offset int) {
	if c.fromRegister != "" {
		m.Write(c.toRegister, m.Read(c.fromRegister))
		return
	}
	m.Write(c.toRegister, c.fromValue)
	return
}

func (c *Copy) String() string {
	if c.fromRegister != "" {
		return fmt.Sprintf("copy(%q, %q)", c.fromRegister, c.toRegister)
	}
	return fmt.Sprintf("copy(%d, %q)", c.fromValue, c.toRegister)
}

type Add struct {
	register string
	d        int
}

func (a *Add) Apply(m Memory) (offset int) {
	m.Write(a.register, m.Read(a.register)+a.d)
	return
}

func (a *Add) String() string {
	return fmt.Sprintf("add(%d, %q)", a.d, a.register)
}

type JumpNonZero struct {
	register string
	offset   int
}

func (j *JumpNonZero) Apply(m Memory) int {
	if j.register != "" && m.Read(j.register) == 0 {
		return 0
	}
	return j.offset - 1
}

func (j *JumpNonZero) String() string {
	if j.register == "" {
		return fmt.Sprintf("jump(%d)", j.offset)
	}
	return fmt.Sprintf("jump(%q, %d)", j.register, j.offset)
}
