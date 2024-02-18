package main

import (
	"chord/chord"
	"fmt"
)

// A small test scenario to test logic between nodes.
// Nodes communicate using the local interface and not RPC
func main() {
	lead := chord.CreateNode(0)
	lead.Start()

	n := 32
	s := make([]*chord.Node, n)

	fmt.Printf("Setting up %v local nodes", n+1)
	for i := 1; i < n; i++ {
		s[i] = chord.CreateNode(chord.Id(1 << i))
		s[i].Join(lead)
		s[i].Start()
	}

	select {}
}
