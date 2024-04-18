package dht

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var promSetKeysTotal = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "dht_set_key_calls_total",
	Help: "Count of SetKey operations",
})

var promGetKeysTotal = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "dht_get_key_calls_total",
	Help: "Count of GetKey operations",
})

var promDeleteKeysTotal = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "dht_delete_keys_total",
	Help: "Count of GetKey operations",
})
