package chord

import (
	"github.com/prometheus/client_golang/prometheus"
)

var successorGauge = prometheus.GaugeOpts{
	Name: "chord_successor",
	Help: "Node's successor",
}

var predecessorGauge = prometheus.GaugeOpts{
	Name: "chord_precessor",
	Help: "Node's successor",
}

var operationsCounter = prometheus.CounterOpts{
	Name: "chord_operation_count_total",
	Help: "Counter of Chord operations",
}

var operationsCounterLabels = []string{"operation", "status", "id"}
