package chord

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var promStabilizeRounds = promauto.NewCounter(prometheus.CounterOpts{
	Name: "chord_stabilize_rounds_total",
	Help: "Count of stabilization rounds started",
})

var promStabilizeRoundsFailed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "chord_stabilize_failed_rounds_total",
	Help: "Count of failed stabilization rounds",
})

var promRectifies = promauto.NewCounter(prometheus.CounterOpts{
	Name: "chord_rectifies_total",
	Help: "Count of successful rectifications",
})

var promFindSuccessors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "chord_find_successor_calls_total",
	Help: "Count of calls to FindSuccessor",
})
