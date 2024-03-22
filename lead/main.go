package main

import (
	"chord/chord"
	"fmt"
	"net"
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
		lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", 8080))
		if err != nil {
			panic("could not start listener")
		}
		chord.StartServer(node, lis)
	}()

	node.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Exiting...")
	node.Stop()
}
