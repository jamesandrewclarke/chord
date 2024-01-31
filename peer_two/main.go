package main

import (
	"chord/chord"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	node := chord.CreateNode(4561)

	chord.SetPeerAddress(4561, "localhost:8080")
	chord.SetPeerAddress(58752, "localhost:8081")
	chord.SetPeerAddress(239847, "localhost:8082")

	go func() {
		chord.StartServer(node, 8080)
	}()

	node.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Exiting...")
}
