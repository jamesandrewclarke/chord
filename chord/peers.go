package chord

import "fmt"

var peers map[Id]node

// getPeer returns a node given an identifier
func getPeer(id Id) (node, error) {
	if peer, ok := peers[id]; ok {
		return peer, nil
	}

	return nil, fmt.Errorf("unable to locate peer %v in cache", id)
}
