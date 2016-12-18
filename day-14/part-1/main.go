package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	salt := os.Args[1]

	N := int(^uint(0) >> 1)
	hashes := make(HashSlice, 0)
	triples := make(map[byte][]Hash)
	for n := 0; n < N; n++ {
		h := hash(salt, n)
		if h.t == 0 {
			continue
		}

		// Check for characters repeating 5 times.
		for c := 0; c < 28; c++ {
			if h.h[c] == h.h[c+1] && h.h[c] == h.h[c+2] && h.h[c] == h.h[c+3] && h.h[c] == h.h[c+4] {
				log.Println("=", h)
				for _, tripHash := range triples[h.h[c]] {
					if tripHash.n+1000 < h.n {
						log.Println("-", tripHash)
						continue
					}

					log.Println("+", tripHash)
					hashes = append(hashes, tripHash)
					if len(hashes) == 64 {
						N = tripHash.n + 1002
					}
				}
				delete(triples, h.h[c])
			}
		}

		// Queue hashes against their triples.
		triples[h.t] = append(triples[h.t], h)
	}

	sort.Sort(hashes)
	for i, h := range hashes {
		fmt.Println(i+1, h)
	}
}

func hash(salt string, index int) Hash {
	h := fmt.Sprintf("%x", md5.Sum([]byte(salt+strconv.Itoa(index))))
	var t byte
	for c := 0; c < 30; c++ {
		if h[c] == h[c+1] && h[c] == h[c+2] {
			t = h[c]
		}
	}
	return Hash{h: h, t: t, n: index}
}

type HashSlice []Hash

func (h HashSlice) Len() int           { return len(h) }
func (h HashSlice) Swap(a, b int)      { h[a], h[b] = h[b], h[a] }
func (h HashSlice) Less(i, j int) bool { return h[i].n < h[j].n }

type Hash struct {
	h string
	t byte
	n int
}

func (h Hash) String() string {
	if h.t != 0 {
		return fmt.Sprintf("%d: %s %s", h.n, h.h, string(h.t))
	}
	return fmt.Sprintf("%d: %s", h.n, h.h)
}
