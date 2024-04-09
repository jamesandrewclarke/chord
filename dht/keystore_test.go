package dht

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyStoreHasKey(t *testing.T) {
	k := CreateKeyStore()
	err := k.SetKey("test", []byte("test"))

	assert.Nil(t, err, "expected nil err")

	assert.True(t, k.HasKey("test"))
}
