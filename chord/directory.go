package chord

import (
	"fmt"
	"sync"
)

type peerStore struct {
	mu            sync.Mutex
	peerAddresses map[Id]node
}

var store peerStore

func init() {
	store.peerAddresses = make(map[Id]node)
}

func SavePeer(node node) {
	if node == nil {
		fmt.Printf("Received nil node for saving")
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if curr, ok := store.peerAddresses[node.Identifier()]; ok && curr != node {
		// log.Printf("overwriting peer for %d\n with: %v", node.Identifier(), node.String())
	}
	store.peerAddresses[node.Identifier()] = node
}

func GetPeer(id Id) (node, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if node, ok := store.peerAddresses[id]; ok {
		return node, nil
	}

	return nil, fmt.Errorf("peer %v not fonud in directory", id)
}
