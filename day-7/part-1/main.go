package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		ip := s.Text()
		if ipSupportsTLS(ip) {
			fmt.Println(ip)
		}
	}

	if err := s.Err(); err != nil {
		panic(err)
	}
}

func ipSupportsTLS(ip string) (support bool) {
	hyper := false
	start := 0

	for c := 0; c < len(ip)+1; c++ {
		if c-start >= 4 {
			if isABBA(ip[c-4 : c]) {
				if hyper {
					return false
				}
				support = true
			}
		}

		if c < len(ip) && (ip[c] == '[' || ip[c] == ']') {
			hyper = ip[c] == '['
			start = c + 1
			continue
		}
	}

	return
}

func isABBA(ip string) bool {
	return ip[0] == ip[3] && ip[1] == ip[2] && ip[0] != ip[1]
}
