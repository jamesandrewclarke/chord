package main

import (
	"chord/chord"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	id := chord.IdentifierFromAddress("127.0.0.1:8080")
	node := chord.CreateNode(id)

	chord.SavePeer(&chord.RPCNode{
		Id:      id,
		Address: "127.0.0.1:8080",
	})

	chord.SetExternalAddress("127.0.0.1:8080")

	go func() {
		chord.StartServer(node, 8080)
	}()

	node.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Exiting...")
	node.Stop()
}
