package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

func init() {
	// Type checks.
	var _ Literal = make(MultiLiteral, 0)
	var _ Literal = &RawLiteral{}
	var _ Literal = &MultiplierLiteral{}
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "-c" {
		count(os.Stdin)
	} else {
		decompress(os.Stdin)
	}
}

func count(r io.Reader) {
	c := &Counter{NewScanner(r)}
	n, err := c.CountRunes()
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}

func decompress(r io.Reader) {
	d := &Decompressor{NewScanner(r)}
	_, err := d.WriteTo(os.Stdout)

	if err != nil {
		panic(err)
	}
}

type Counter struct {
	Scanner *Scanner
}

func (c *Counter) CountRunes() (int64, error) {
	var n int64
	for {
		_, tok, lit, err := c.Scanner.Scan(0)
		if err != nil {
			return 0, err
		}
		if tok == TokenEOF {
			break
		}
		n += lit.Length()
	}
	return n, nil
}

type Decompressor struct {
	Scanner *Scanner
}

func (d *Decompressor) WriteTo(w io.Writer) (int64, error) {
	var n int64
	for {
		_, tok, lit, err := d.Scanner.Scan(0)
		if err != nil {
			return n, err
		}
		if tok == TokenEOF {
			break
		}

		nn, err := lit.WriteTo(w)
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Scan(limit int) (int, Token, Literal, error) {
	fmt.Fprintf(os.Stderr, "Scan(%d)\n", limit)
	ch := s.read()

	if isOpenParen(ch) {
		s.unread()
		return s.scanMultiplier()
	} else if isEOF(ch) {
		return 0, TokenEOF, &RawLiteral{}, nil
	}

	s.unread()
	return s.scanRaw(limit)
}

func (s *Scanner) ScanMulti(limit int) (int, Token, Literal, error) {
	fmt.Fprintf(os.Stderr, "ScanMulti(%d)\n", limit)
	n := 0
	token := TokenMulti
	literal := make(MultiLiteral, 0)

	all := limit < 1
	for all || limit >= 1 {
		fmt.Fprintf(os.Stderr, "+ScanMulti(%d)\n", limit)
		dn, t, l, err := s.Scan(limit)
		if t == TokenEOF {
			break
		}

		n += dn
		limit -= dn
		if err != nil {
			return n, t, l, err
		}
		literal = append(literal, l)
	}

	return n, token, literal, nil
}

func (s *Scanner) scanRaw(limit int) (int, Token, Literal, error) {
	fmt.Fprintf(os.Stderr, "scanRaw(%d)\n", limit)
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var n int
	for n = 1; limit < 1 || n < limit; {
		ch := s.read()
		if isOpenParen(ch) || isEOF(ch) {
			s.unread()
			break
		}
		n++
		buf.WriteRune(ch)
	}

	return n, TokenString, &RawLiteral{buf.Bytes()}, nil
}

func (s *Scanner) scanMultiplier() (int, Token, Literal, error) {
	fmt.Fprintf(os.Stderr, "scanMultiplier()\n")
	buf := new(bytes.Buffer)
	n := 5 // The number of tokens read outside of loops.

	ch := s.read()
	if !isOpenParen(ch) {
		return 1, TokenIllegal, nil, fmt.Errorf("found %q, expected '('", ch)
	}

	ch = s.read()
	if !isDigit(ch) {
		return 2, TokenIllegal, nil, fmt.Errorf("found %q, expected digit", ch)
	}
	buf.WriteRune(ch)

	for {
		ch := s.read()
		if ch == 'x' {
			break
		} else if !isDigit(ch) {
			return n - 2, TokenIllegal, nil, fmt.Errorf("found %q, expected digit or 'x'", ch)
		}
		n++
		buf.WriteRune(ch)
	}

	length, _ := strconv.Atoi(string(buf.Bytes()))
	buf.Reset()

	ch = s.read()
	if !isDigit(ch) {
		return n - 1, TokenIllegal, nil, fmt.Errorf("found %q, expected digit", ch)
	}
	buf.WriteRune(ch)

	for {
		ch := s.read()
		if isCloseParen(ch) {
			break
		} else if !isDigit(ch) {
			return n, TokenIllegal, nil, fmt.Errorf("found %q, expected digit or ')'", ch)
		}
		n++
		buf.WriteRune(ch)
	}

	times, _ := strconv.Atoi(string(buf.Bytes()))

	innerN, _, repeatedLiteral, err := s.ScanMulti(length)
	if err != nil {
		return n + innerN, TokenIllegal, nil, fmt.Errorf("failed to read repeated token: %s", err)
	}
	if innerN != length {
		return n + innerN, TokenIllegal, nil, fmt.Errorf("tried to read %d runes but found %d: %+v", length, innerN, repeatedLiteral)
	}
	return n + length, TokenString, &MultiplierLiteral{times, repeatedLiteral}, nil
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		fmt.Fprintln(os.Stderr, "read: eof")
		return eof
	}
	fmt.Fprintf(os.Stderr, "read: %q\n", ch)
	return ch
}

func (s *Scanner) unread() {
	fmt.Fprintf(os.Stderr, "unread\n")
	s.r.UnreadRune()
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
func isOpenParen(ch rune) bool {
	return ch == '('
}

func isCloseParen(ch rune) bool {
	return ch == ')'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isEOF(ch rune) bool {
	return ch == eof
}

var eof = rune(0)

type Token int

const (
	TokenIllegal Token = iota
	TokenEOF
	TokenString
	TokenParenOpen
	TokenParenClose
	TokenInteger
	TokenMulti
)

type Literal interface {
	io.WriterTo
	Length() int64
}

type MultiLiteral []Literal

func (m MultiLiteral) WriteTo(w io.Writer) (int64, error) {
	var n int64
	for _, l := range m {
		dn, err := l.WriteTo(w)
		n += dn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (m MultiLiteral) Length() (n int64) {
	for _, l := range m {
		n += l.Length()
	}
	return
}

type RawLiteral struct {
	body []byte
}

func (r *RawLiteral) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(r.body)
	return int64(n), err
}

func (r *RawLiteral) Length() int64 {
	return int64(len(r.body))
}

func (r *RawLiteral) Body() []byte {
	return r.body
}

type MultiplierLiteral struct {
	times           int
	repeatedLiteral Literal
}

func (m *MultiplierLiteral) WriteTo(w io.Writer) (int64, error) {
	var n int64
	for i := 0; i < m.times; i++ {
		nn, err := m.repeatedLiteral.WriteTo(w)
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (m *MultiplierLiteral) Length() (n int64) {
	return int64(m.times) * m.repeatedLiteral.Length()
}
