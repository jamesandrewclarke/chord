package chord

import (
	"fmt"
	"log/slog"
	"sync"
)

type keystore interface {
	HasKey(Id) bool
	SetKey(Id, []byte) error
	GetKey(Id) ([]byte, error)
	DeleteKey(Id) error
}

type KeyStore struct {
	muKeys sync.RWMutex
	keys   map[Id][]byte
}

func CreateKeyStore() *KeyStore {
	ks := &KeyStore{}
	ks.keys = make(map[Id][]byte)

	return ks
}

func (k *KeyStore) HasKey(id Id) bool {
	k.muKeys.RLock()
	defer k.muKeys.RUnlock()

	_, ok := k.keys[id]
	return ok
}

func (k *KeyStore) SetKey(id Id, bytes []byte) error {
	k.muKeys.Lock()
	defer k.muKeys.Unlock()

	if k.hasKey(id) {
		slog.Warn("overwriting log entry", "id", id)
	}

	k.keys[id] = bytes
	return nil
}

// hasKey is the non-threadsafe version of HasKey for internal use only
func (k *KeyStore) hasKey(id Id) bool {
	_, ok := k.keys[id]
	return ok
}

func (k *KeyStore) GetKey(id Id) ([]byte, error) {
	k.muKeys.RLock()
	defer k.muKeys.RUnlock()

	if !k.hasKey(id) {
		return nil, fmt.Errorf("key %v not found", id)
	}

	return k.keys[id], nil
}

func (k *KeyStore) DeleteKey(id Id) error {
	k.muKeys.Lock()
	defer k.muKeys.Unlock()

	if !k.hasKey(id) {
		return fmt.Errorf("could not delete key %v: not found", id)
	}

	delete(k.keys, id)
	return nil
}
