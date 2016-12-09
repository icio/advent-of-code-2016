package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println(doorPassword(os.Args[1]))
}

func doorPassword(doorId string) string {
	pass := make([]byte, 8)
	var index int

	for count := 0; count < 8; {
		index++

		hash := fmt.Sprintf("%x", md5.Sum([]byte(doorId+strconv.Itoa(index))))
		if hash[0:5] != "00000" {
			continue
		}

		pos := hash[5] - 48 // ord('0') == 48
		if pos < 0 || pos > 7 || pass[pos] != 0 {
			continue
		}

		count += 1
		pass[pos] = hash[6]
		fmt.Fprintln(os.Stderr, hash, pos, count, index, string(bytes.Replace(pass, []byte{0}, []byte{'_'}, -1)))
	}

	return string(pass)
}
