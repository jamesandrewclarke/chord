package chord

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"net"
)

var HashFunc = sha256.Sum256

type BootstrapConfig struct {
	// ExternalAddr is the address that peers will contact the node on, this
	// ultimately determines the node's identifier. Should not include the port.
	ExternalAddr string

	// The port to reach the node on, if unspecified, a random port will be chosen
	Port int

	// An address and port of an existing node in the desired network, if unspecified,
	// the ring will initialise with the single new node
	BootstrapAddr string
}

func Bootstrap(config BootstrapConfig) *Node {
	lis := getListener(config.Port)
	port := lis.Addr().(*net.TCPAddr).Port

	addr := fmt.Sprintf("%v:%v", config.ExternalAddr, port)
	SetExternalAddress(addr)

	var node *Node
	if config.BootstrapAddr != "" {
		lead_id := IdentifierFromAddress(config.BootstrapAddr)
		remote := &RPCNode{
			Address: config.BootstrapAddr,
			Id:      lead_id,
		}
		SavePeer(remote)

		id := remote.Announce(port, nil)
		node = CreateNode(id)
		SavePeer(node)

		err := node.Join(remote)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("No bootstrap address provided, initialising a new Chord ring")
		id := IdentifierFromAddress(addr)
		node = CreateNode(id)
		SavePeer(node)
	}

	node.Start()

	go func() {
		StartServer(node, lis)
	}()

	return node
}

func getListener(port int) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		panic("could not start listener")
	}

	return lis
}

// IdentifierFromAddress takes a peer address and computes its identifier using a hash function
// addr should be of the form <ip address>:port
func IdentifierFromAddress(addr string) Id {
	sum := HashFunc([]byte(addr))

	return IdentifierFromBytes(sum[:])
}

// Hash returns a slice of the checksum calculated using HashFunc
func Hash(bytes []byte) []byte {
	sum := HashFunc(bytes)
	return sum[:]
}

// Produce an Id modulo m given some arbitrary bytes
func IdentifierFromBytes(bytes []byte) Id {
	bigint := new(big.Int)
	bigint.SetBytes(bytes)

	bigid := new(big.Int)

	bigid.Mod(bigint, big.NewInt(1<<(m-1)-1))
	// bigid.Mod(bigint, big.NewInt(1<<11))

	return Id(bigid.Int64())
}

// Between returns true if id is in the range [start, end] on the Chord ring
func Between(id, start, end Id) bool {
	if start < end {
		return id > start && id < end
	}

	return id > start || id < end
}

// NodesBetween returns if Id falls between the identifiers of nodes a and b (exclusive)
func NodesBetween(id Id, a, b node) bool {
	return Between(a.Identifier(), b.Identifier(), id)
}
