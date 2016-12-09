package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		ip := s.Text()
		if ipSupportsSSL(ip) {
			fmt.Println(ip)
		}
	}

	if err := s.Err(); err != nil {
		panic(err)
	}
}

func ipSupportsSSL(ip string) bool {
	base, hyper := "", ""

	for s := 0; ; {
		c := strings.IndexAny(ip[s:], "[]")
		if c == -1 {
			base += " " + ip[s:]
			break
		} else if ip[s+c] == '[' {
			base += " " + ip[s:s+c]
		} else if ip[s+c] == ']' {
			hyper += " " + ip[s:s+c]
		}
		s += c + 1
	}

	for s := 1; s < len(base)-2; s++ {
		if base[s+2] == ' ' {
			s += 2
			continue
		}

		aba := base[s : s+3]
		if aba[0] != aba[2] || aba[0] == aba[1] {
			continue
		}

		if strings.Index(hyper, string([]byte{aba[1], aba[0], aba[1]})) > -1 {
			return true
		}
	}

	return false
}
