package main

import (
	"chord/chord"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var BOOTSTRAP_ADDRESS = flag.String("addr", "", "The address and port of an existing node to join")

var PORT = flag.Int("port", 0, "Port to listen on")

// getListener returns a listener bound to 0.0.0.0 on a random port (unless specified)
func getListener() net.Listener {
	// if *PORT is 0, the net library will just assign a port itself
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *PORT))
	if err != nil {
		panic("could not start listener")
	}

	return lis
}

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

	lis := getListener()
	port := lis.Addr().(*net.TCPAddr).Port
	chord.SetExternalAddress(fmt.Sprintf("127.0.0.1:%v", port))

	id := remote.Announce(port, nil)
	node := chord.CreateNode(id)

	chord.SavePeer(node)
	chord.SavePeer(remote)

	err := node.Join(remote)
	if err != nil {
		panic(err)
	}

	node.Start()

	go func() {
		chord.StartServer(node, lis)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Exiting...")
	node.Stop()
}
