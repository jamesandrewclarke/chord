package dht

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyStoreHasKey(t *testing.T) {
	k := CreateKeyStore(0)
	err := k.SetKey("test", []byte("test"))

	assert.Nil(t, err, "expected nil err")

	assert.True(t, k.HasKey("test"))
}

func TestKeyStoreReturnsCorrectKey(t *testing.T) {
	k := CreateKeyStore(0)

	err := k.SetKey("test", []byte("Hello, World!"))

	assert.Nil(t, err, "expected nil err")

	value, err := k.GetKey("test")
	assert.Nil(t, err, "expected nil err")
	assert.Equal(t, []byte("Hello, World!"), value)
}
