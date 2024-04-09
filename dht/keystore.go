package dht

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var promKeysTotal = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "dht_keys_total",
	Help: "The total number of keys stored in the node",
})

type keystore interface {
	HasKey(string) bool
	SetKey(string, []byte) error
	GetKey(string) ([]byte, error)
	DeleteKey(string) error
}

type keyentry struct {
	key   string
	Value []byte

	sync.RWMutex
}

type KeyStore struct {
	muKeys sync.RWMutex
	keys   map[string]*keyentry
}

func CreateKeyStore() *KeyStore {
	ks := &KeyStore{}
	ks.keys = make(map[string]*keyentry)

	return ks
}

func createKeyEntry(key string) *keyentry {
	k := &keyentry{}
	k.key = key
	return k
}

func (k *KeyStore) HasKey(key string) bool {
	k.muKeys.RLock()
	defer k.muKeys.RUnlock()

	_, ok := k.keys[key]
	return ok
}

func (k *KeyStore) SetKey(key string, bytes []byte) error {
	k.muKeys.Lock()
	defer k.muKeys.Unlock()

	if k.hasKey(key) {
		slog.Warn("overwriting log entry", "key", key)
	} else {
		promKeysTotal.Inc()
		k.keys[key] = createKeyEntry(key)
	}

	entry := k.keys[key]
	entry.Lock()
	defer entry.Unlock()
	entry.Value = bytes

	return nil
}

// hasKey is the non-threadsafe version of HasKey for internal use only
func (k *KeyStore) hasKey(key string) bool {
	_, ok := k.keys[key]
	return ok
}

func (k *KeyStore) GetKey(key string) ([]byte, error) {
	k.muKeys.RLock()
	defer k.muKeys.RUnlock()

	if !k.hasKey(key) {
		return nil, fmt.Errorf("key %v not found", key)
	}
	entry := k.keys[key]
	entry.Lock()
	defer entry.Unlock()

	return entry.Value, nil
}

func (k *KeyStore) DeleteKey(key string) error {
	k.muKeys.Lock()
	defer k.muKeys.Unlock()

	if !k.hasKey(key) {
		return fmt.Errorf("could not delete key %v: not found", key)
	}

	entry := k.keys[key]
	entry.Lock()
	defer entry.Unlock()
	delete(k.keys, key)

	promKeysTotal.Dec()
	return nil
}
