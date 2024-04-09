package dht

import "net"

func stripPort(address string) string {
	host, _, _ := net.SplitHostPort(address)
	return host
}
