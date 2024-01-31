package chord

import "fmt"

var peers map[Id]node

var peerAddresses map[Id]string

func init() {
	peerAddresses = make(map[Id]string)
}

// getPeer returns a node given an identifier
func getPeer(id Id) (node, error) {
	if peer, ok := peers[id]; ok {
		return peer, nil
	}

	return nil, fmt.Errorf("unable to locate peer %v in cache", id)
}

func SetPeerAddress(id Id, addr string) {
	peerAddresses[id] = addr
}

func getPeerAddress(id Id) (string, error) {
	if addr, ok := peerAddresses[id]; ok {
		return addr, nil
	}

	return "", fmt.Errorf("unable to locate peer %v in cache", id)
}
