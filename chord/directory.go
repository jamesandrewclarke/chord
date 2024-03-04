package chord

import "fmt"

var peerAddresses map[Id]string

func init() {
	peerAddresses = make(map[Id]string)
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
