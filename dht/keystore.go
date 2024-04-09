package dht

import (
	"chord_dht/chord"
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
	HasKey(chord.Id) bool
	SetKey(chord.Id, []byte) error
	GetKey(chord.Id) ([]byte, error)
	DeleteKey(chord.Id) error
}

type KeyStore struct {
	muKeys sync.RWMutex
	keys   map[chord.Id][]byte
}

func CreateKeyStore() *KeyStore {
	ks := &KeyStore{}
	ks.keys = make(map[chord.Id][]byte)

	return ks
}

func (k *KeyStore) HasKey(id chord.Id) bool {
	k.muKeys.RLock()
	defer k.muKeys.RUnlock()

	_, ok := k.keys[id]
	return ok
}

func (k *KeyStore) SetKey(id chord.Id, bytes []byte) error {
	k.muKeys.Lock()
	defer k.muKeys.Unlock()

	if k.hasKey(id) {
		slog.Warn("overwriting log entry", "id", id)
	} else {
		promKeysTotal.Inc()
	}

	k.keys[id] = bytes
	return nil
}

// hasKey is the non-threadsafe version of HasKey for internal use only
func (k *KeyStore) hasKey(id chord.Id) bool {
	_, ok := k.keys[id]
	return ok
}

func (k *KeyStore) GetKey(id chord.Id) ([]byte, error) {
	k.muKeys.RLock()
	defer k.muKeys.RUnlock()

	if !k.hasKey(id) {
		return nil, fmt.Errorf("key %v not found", id)
	}

	return k.keys[id], nil
}

func (k *KeyStore) DeleteKey(id chord.Id) error {
	k.muKeys.Lock()
	defer k.muKeys.Unlock()

	if !k.hasKey(id) {
		return fmt.Errorf("could not delete key %v: not found", id)
	}

	promKeysTotal.Dec()
	delete(k.keys, id)
	return nil
}
