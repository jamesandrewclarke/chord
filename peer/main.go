package main

import (
	"chord/chord"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var BOOTSTRAP_ADDRESS = flag.String("addr", "", "The address and port of an existing node to join")

var PORT = flag.Int("port", 8081, "Port to listen on")

func main() {
	flag.Parse()

	lead_id := chord.IdentifierFromAddress(*BOOTSTRAP_ADDRESS)
	remote := &chord.RPCNode{
		Address: *BOOTSTRAP_ADDRESS,
		Id:      lead_id,
	}

	if BOOTSTRAP_ADDRESS == nil {
		panic("no bootstrap address provided")
	}

	id := remote.Announce(*PORT, nil)
	node := chord.CreateNode(id)

	// yes
	chord.SetPeerAddress(id, fmt.Sprintf("127.0.0.1:%v", *PORT))
	chord.SetPeerAddress(lead_id, *BOOTSTRAP_ADDRESS)

	err := node.Join(remote)
	if err != nil {
		panic(err)
	}

	node.Start()

	go func() {
		chord.StartServer(node, *PORT)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Exiting...")
	node.Stop()
}
