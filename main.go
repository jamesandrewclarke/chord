package main

import (
	"chord/chord"
)

func main() {
	lead := chord.CreateNode(1)
	lead.Start()

	s := make([]*chord.Node, 10)
	for i := 1; i < 10; i++ {
		s[i] = chord.CreateNode(chord.Id(i) + 1)
		s[i].Join(lead)
		s[i].Start()
	}

	for {
	}
}
