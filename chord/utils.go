package chord

import (
	"crypto/sha256"
	"math/big"
)

// IdentifierFromAddress takes a peer address and computes its identifier using a hash function
// addr should be of the form <ip address>:port
func IdentifierFromAddress(addr string) Id {
	sum := sha256.Sum256([]byte(addr))

	return IdentifierFromBytes(sum[:])
}

func IdentifierFromBytes(bytes []byte) Id {
	bigint := new(big.Int)
	bigint.SetBytes(bytes)

	bigid := new(big.Int)

	bigid.Mod(bigint, big.NewInt(1<<(m-1)-1))
	// bigid.Mod(bigint, big.NewInt(1<<11))

	return Id(bigid.Int64())
}
