package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	d := &Decompressor{NewScanner(os.Stdin)}
	_, err := d.WriteTo(os.Stdout)

	if err != nil {
		panic(err)
	}
}

type Decompressor struct {
	Scanner *Scanner
}

func (d *Decompressor) WriteTo(w io.Writer) (int64, error) {
	var n int64
	for {
		tok, lit, err := d.Scanner.Scan()
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

func (s *Scanner) Scan() (Token, Literal, error) {
	ch := s.read()

	if isOpenParen(ch) {
		return s.scanMultiplier()
	} else if isEOF(ch) {
		return TokenEOF, &RawLiteral{}, nil
	}

	s.unread()
	return s.scanRaw()
}

func (s *Scanner) scanRaw() (Token, Literal, error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		if isOpenParen(ch) || isEOF(ch) {
			s.unread()
			break
		}
		buf.WriteRune(ch)
	}

	return TokenString, &RawLiteral{buf.Bytes()}, nil
}

func (s *Scanner) scanMultiplier() (Token, Literal, error) {
	buf := new(bytes.Buffer)

	ch := s.read()
	if !isDigit(ch) {
		return TokenIllegal, nil, fmt.Errorf("found %q, expected digit", ch)
	}
	buf.WriteRune(ch)

	for {
		ch := s.read()
		if ch == 'x' {
			break
		} else if !isDigit(ch) {
			return TokenIllegal, nil, fmt.Errorf("found %q, expected digit or 'x'", ch)
		}
		buf.WriteRune(ch)
	}

	length, _ := strconv.Atoi(string(buf.Bytes()))
	buf.Reset()

	ch = s.read()
	if !isDigit(ch) {
		return TokenIllegal, nil, fmt.Errorf("found %q, expected digit", ch)
	}
	buf.WriteRune(ch)

	for {
		ch := s.read()
		if isCloseParen(ch) {
			break
		} else if !isDigit(ch) {
			return TokenIllegal, nil, fmt.Errorf("found %q, expected digit or ')'", ch)
		}
		buf.WriteRune(ch)
	}

	times, _ := strconv.Atoi(string(buf.Bytes()))

	buf.Reset()
	n, err := io.Copy(buf, io.LimitReader(s.r, int64(length)))

	if n < int64(length) {
		return TokenIllegal, nil, fmt.Errorf("failed to read %s chars: only read %s", length, n)
	}
	if err != nil {
		return TokenIllegal, nil, fmt.Errorf("failed to read %d chars: %s", length, err)
	}

	return TokenString, &MultiplierLiteral{times, buf.Bytes()}, nil
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
	fmt.Fprintf(os.Stderr, "unread")
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
)

type Literal interface {
	io.WriterTo
	Body() []byte
}

type RawLiteral struct {
	body []byte
}

func (r *RawLiteral) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(r.body)
	return int64(n), err
}

func (r *RawLiteral) Body() []byte {
	return r.body
}

type MultiplierLiteral struct {
	times int
	body  []byte
}

func (m *MultiplierLiteral) WriteTo(w io.Writer) (int64, error) {
	var n int64
	for i := 0; i < m.times; i++ {
		nn, err := w.Write(m.body)
		n += int64(nn)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (m *MultiplierLiteral) Body() []byte {
	return m.body
}
