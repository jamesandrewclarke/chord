package main

import (
	"chord_dht/chord"
	"chord_dht/dht"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var EXTERNAL_ADDRESS = flag.String("address", "127.0.0.1", "The address that peers will contact the server on, should be set accordingly for networks behind a NAT")

var BOOTSTRAP_ADDRESS = flag.String("bootstrap", "", "The address and port of a node in an existing Chord ring")

var PORT = flag.Int("port", 0, "Port to listen on")

func main() {
	flag.Parse()

	config := chord.BootstrapConfig{
		ExternalAddr:  *EXTERNAL_ADDRESS,
		BootstrapAddr: *BOOTSTRAP_ADDRESS,
		Port:          *PORT,
	}
	node := chord.Bootstrap(config)

	go func() {
		dht.StartDHT(node, 8081)
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Exiting...")
	node.Stop()
}
