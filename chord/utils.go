package chord

import (
	"crypto/sha256"
	"math/big"
)

var HashFunc = sha256.Sum256

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
