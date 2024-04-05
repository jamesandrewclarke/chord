package main

import (
	"chord/chord"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var EXTERNAL_ADDRESS = flag.String("address", "127.0.0.1", "The address that peers will contact the server on, should be set accordingly for networks behind a NAT")

var BOOTSTRAP_ADDRESS = flag.String("bootstrap", "", "The address and port of a node in an existing Chord ring")

// The default of 0 ensures a random port is assigned by the OS
var PORT = flag.Int("port", 0, "Port to listen on")

func getListener() net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *PORT))
	if err != nil {
		panic("could not start listener")
	}

	return lis
}

func main() {
	flag.Parse()

	lis := getListener()
	port := lis.Addr().(*net.TCPAddr).Port

	addr := fmt.Sprintf("%v:%v", *EXTERNAL_ADDRESS, port)

	fmt.Printf("Chord address: %v\n", addr)
	chord.SetExternalAddress(addr)

	var node *chord.Node
	if *BOOTSTRAP_ADDRESS != "" {
		lead_id := chord.IdentifierFromAddress(*BOOTSTRAP_ADDRESS)
		remote := &chord.RPCNode{
			Address: *BOOTSTRAP_ADDRESS,
			Id:      lead_id,
		}
		chord.SavePeer(remote)

		id := remote.Announce(port, nil)
		node = chord.CreateNode(id)
		chord.SavePeer(node)

		err := node.Join(remote)
		if err != nil {
			panic(err)
		}
	} else {
		id := chord.IdentifierFromAddress(addr)
		node = chord.CreateNode(id)
		chord.SavePeer(node)
	}

	node.Start()

	go func() {
		chord.StartServer(node, lis)
	}()

	go func() {
		// Prometheus for instrumentation
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Exiting...")
	node.Stop()
}
