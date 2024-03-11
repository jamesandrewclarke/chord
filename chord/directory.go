package chord

import (
	"fmt"
	"log"
	"net"
	"sync"
)

var mu sync.Mutex

var peerAddresses map[Id]string

func init() {
	peerAddresses = make(map[Id]string)
}

func SetPeerAddress(id Id, addr string) {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Printf("Invalid address for %v: `%v`", id, addr)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if currAddr, ok := peerAddresses[id]; ok && currAddr != addr {
		log.Printf("overwriting address for %d\n from %v to %v", id, currAddr, addr)
	}
	peerAddresses[id] = addr
}

func getPeerAddress(id Id) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	if addr, ok := peerAddresses[id]; ok {
		return addr, nil
	}

	return "", fmt.Errorf("unable to locate peer %v in cache", id)
}
