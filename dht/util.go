package dht

import (
	"chord_dht/chord"
	"net"
)

func stripPort(address string) string {
	host, _, _ := net.SplitHostPort(address)
	return host
}

func ChordIdFromString(str string) chord.Id {
	hash := chord.Hash([]byte(str))
	return chord.IdentifierFromBytes(hash)
}
