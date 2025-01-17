package dht

import (
	"chord_dht/chord"
	"fmt"
	"log/slog"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type keystore interface {
	HasKey(string) bool
	SetKey(string, []byte) error
	GetKey(string) ([]byte, error)
	DeleteKey(string) error
}

type keyentry struct {
	Key   string
	Value []byte
	Id    chord.Id

	sync.RWMutex
}

type KeyStore struct {
	muKeys sync.RWMutex
	Keys   map[string]*keyentry

	keyGauge prometheus.Gauge

	Registry *prometheus.Registry
}

func CreateKeyStore(id chord.Id) *KeyStore {
	ks := &KeyStore{
		Keys: make(map[string]*keyentry),
		keyGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "dht_keys_total",
			Help: "The total number of keys stored in the node",
			ConstLabels: prometheus.Labels{
				"id": fmt.Sprint(id),
			},
		}),
	}

	ks.Registry = prometheus.NewRegistry()
	ks.Registry.MustRegister(ks.keyGauge)

	return ks
}

func createKeyEntry(key string) *keyentry {
	k := &keyentry{}
	k.Key = key
	k.Id = ChordIdFromString(key)
	return k
}

func (k *KeyStore) HasKey(key string) bool {
	k.muKeys.RLock()
	defer k.muKeys.RUnlock()

	_, ok := k.Keys[key]
	return ok
}

func (k *KeyStore) SetKey(key string, bytes []byte) error {
	k.muKeys.Lock()

	defer k.muKeys.Unlock()

	if k.hasKey(key) {
		slog.Warn("overwriting log entry", "key", key)
	} else {
		k.keyGauge.Inc()
		k.Keys[key] = createKeyEntry(key)
	}

	entry := k.Keys[key]

	entry.Lock()
	defer entry.Unlock()
	entry.Value = bytes

	promSetKeysTotal.Inc()
	return nil
}

// hasKey is the non-threadsafe version of HasKey for internal use only
func (k *KeyStore) hasKey(key string) bool {
	_, ok := k.Keys[key]
	return ok
}

func (k *KeyStore) GetKey(key string) ([]byte, error) {
	k.muKeys.RLock()
	defer k.muKeys.RUnlock()

	if !k.hasKey(key) {
		return nil, fmt.Errorf("key %v not found", key)
	}
	entry := k.Keys[key]
	entry.Lock()
	defer entry.Unlock()

	promGetKeysTotal.Inc()
	return entry.Value, nil
}

func (k *KeyStore) DeleteKey(key string) error {
	k.muKeys.Lock()
	defer k.muKeys.Unlock()

	if !k.hasKey(key) {
		return fmt.Errorf("could not delete key %v: not found", key)
	}

	entry := k.Keys[key]
	entry.Lock()
	defer entry.Unlock()
	delete(k.Keys, key)

	k.keyGauge.Dec()
	promDeleteKeysTotal.Inc()
	return nil
}
