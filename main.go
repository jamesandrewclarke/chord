package main

import (
	"chord/chord"
	"fmt"
)

func main() {
	N1 := chord.CreateNode(1)
	N10 := chord.CreateNode(10)

	N10.Join(N1)

	N10.Stabilize()
	N1.Stabilize()

	fmt.Println(N1)
	fmt.Println(N10)
}
