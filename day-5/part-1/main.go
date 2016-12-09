package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println(doorPassword(os.Args[1]))
}

func doorPassword(doorId string) (password string) {
	var index int
	for i := 0; i < 8; i++ {
		var hash string
		for {
			hash = fmt.Sprintf("%x", md5.Sum([]byte(doorId+strconv.Itoa(index))))
			if hash[0:5] != "00000" {
				index++
				continue
			}

			password = password + string(hash[5])
			fmt.Fprintln(os.Stderr, i+1, hash, index, password)
			index++
			break
		}
	}

	return
}
