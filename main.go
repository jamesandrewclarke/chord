package main

import (
	"chord/chord"
	"fmt"
)

func main() {
	n := 1024
	s := make([]*chord.Node, n)

	s[0] = chord.CreateNode(0)
	s[0].Start()

	for i := 1; i < n; i++ {
		// s[i] = chord.CreateNode(chord.Id(1 << i))
		s[i] = chord.CreateNode(chord.Id(i * 16))
		s[i].Join(s[i-1])
		s[i].Start()
	}

	fmt.Println("set up...starting")

	for {
	}
}
