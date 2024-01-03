package main

import (
	"chord/chord"
	"fmt"
)

func main() {
	lead := chord.CreateNode(1)
	lead.Start()

	n := 32
	s := make([]*chord.Node, n)
	for i := 1; i < n; i++ {
		s[i] = chord.CreateNode(chord.Id(1 << i))
		s[i].Join(lead)
		s[i].Start()
	}

	fmt.Println("set up...starting")

	for {
	}
}
