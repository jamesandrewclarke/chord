package chord

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Type alias for Chord IDs
type Id int64

// m is the size of the Chord ring, i.e. the ring is modulo (1<<m)
const m = 32

// The stabilization interval in milliseconds
const STABILIZE_INTERVAL = 1000 * time.Millisecond
const FINGER_INTERVAL = 500 * time.Millisecond

// Abstract interface for a node
type node interface {
	Identifier() Id
	Successor() (node, error)
	Predecessor() (node, error)
	FindSuccessor(Id, int) (node, int, error)
	Rectify(node) error
	SuccessorList() (*SuccessorList, error)
	Alive() bool
	String() string
}

type Node struct {
	id Id

	muPred      sync.Mutex
	predecessor node

	muFinger sync.Mutex
	finger   [m]node

	successorList *SuccessorList

	nextFinger int

	shutdown chan struct{}
	wg       *sync.WaitGroup

	// Metrics
	registry         *prometheus.Registry
	operationCount   *prometheus.CounterVec
	successorGauge   prometheus.Gauge
	predecessorGauge prometheus.Gauge
}

type ChordConfig struct {
	// SuccessorListLength is the number of nodes ahead in the ring to store
	// as potential successors
	SuccessorListLength int

	// StabilizeInterval is the number of milliseconds between stabilize operations
	StabilizeInterval int

	// FingerInterval is the number of milliseconds between finger checks
	FingerInterval int

	// InvariantMonitoring controls whether local invariant checking is run
	InvariantMonitoring bool
}

// CreateNode initialises a single-node Chord ring
func CreateNode(Id Id) *Node {
	n := &Node{
		id:            Id,
		shutdown:      make(chan struct{}),
		wg:            new(sync.WaitGroup),
		successorList: CreateSuccessorList(10),

		operationCount: prometheus.NewCounterVec(operationsCounter, operationsCounterLabels),

		successorGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "chord_successor",
			Help: "The successor of the node",
			ConstLabels: prometheus.Labels{
				"id": fmt.Sprint(Id),
			},
		}),
		predecessorGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "chord_predecessor",
			Help: "The predecessor of the node",
			ConstLabels: prometheus.Labels{
				"id": fmt.Sprint(Id),
			},
		}),
	}

	n.setSuccessor(n)
	n.predecessor = n
	n.nextFinger = 1

	n.registry = prometheus.NewRegistry()

	n.registry.MustRegister(n.operationCount)
	n.registry.MustRegister(n.successorGauge)
	n.registry.MustRegister(n.predecessorGauge)

	return n
}

// Identifier returns the m-bit identifier which determines the node's location on the ring
func (n *Node) Identifier() Id {
	return n.id
}

// Predecessor returns a pointer to n's predecessor
func (n *Node) Predecessor() (node, error) {
	if n.predecessor == nil {
		return nil, fmt.Errorf("no known predecessor")
	}
	return n.predecessor, nil
}

// Successor returns a pointer to n's successor
func (n *Node) Successor() (node, error) {
	return n.successorList.Head(), nil
}

// Start starts the background tasks to stabilize n's pointers and lookup table
func (n *Node) Start() {
	n.wg.Add(1)
	go func() {
		stabilizeTicker := time.NewTicker(STABILIZE_INTERVAL)
		defer stabilizeTicker.Stop()

		fingerTicker := time.NewTicker(FINGER_INTERVAL)
		defer fingerTicker.Stop()
		defer n.wg.Done()

		// TODO Configurable intervals for experiments
		for {
			select {
			case <-stabilizeTicker.C:
				n.checkPredecessor()

				err := n.stabilize()
				if err != nil {
					slog.Error("failed stabilization", "node", n.Identifier(), "err", err)
					n.operationCount.WithLabelValues("stabilize", "fail", fmt.Sprint(n.Identifier())).Inc()
				}
				if !n.successorList.UniqueSuccessors() {
					// slog.Warn("duplicate successors", "node", n.Identifier())
				}

				if !n.successorList.Ordered() {
					// slog.Warn("disordered successors", "node", n.Identifier())
				}

				n.operationCount.WithLabelValues("stabilize", "success", fmt.Sprint(n.Identifier())).Inc()
				slog.Debug("successor list", "list", n.successorList.String())

			case <-fingerTicker.C:
				n.fixFingers()

			case <-n.shutdown:
				return
			}
		}
	}()
}

func (n *Node) Stop() {
	close(n.shutdown)
	slog.Info("graceful shutdown")
	n.wg.Wait()
	slog.Info("shutdown complete")
}

// Join joins a Chord ring containing the node p
func (n *Node) Join(p node) error {
	n.predecessor = nil

	succ, _, err := p.FindSuccessor(n.Identifier(), 0)
	if err != nil {
		return err
	}
	n.setSuccessor(succ)

	return nil
}

// setSuccessor is a safe wrapper method for setting n's immediate successor
func (n *Node) setSuccessor(p node) {
	if p == nil {
		return
	}
	n.successorList.SetHead(p)

	// TODO do we need separate locations?
	n.muFinger.Lock()
	defer n.muFinger.Unlock()
	n.finger[0] = p

	n.successorGauge.Set(float64(n.successorList.Head().Identifier()))
}

// stabilize updates the successor list and informs the immediate successor of the node's presence
func (n *Node) stabilize() error {
	defer func() {
		succ, _ := n.Successor()
		if succ != nil {
			err := succ.Rectify(n)
			if err != nil {
				slog.Error("failed rectify", "err", err, "successor", succ)
			}
		}
	}()

	var succStart node
	succ, err := n.Successor()
	if err != nil || succ == nil {
		return fmt.Errorf("can't stabilize, no successor")
	}
	succStart = succ

	succ_pred, err := succ.Predecessor()
	if err != nil {
		// Assume successor to be dead
		n.successorList.PopHead()
		n.setSuccessor(n.successorList.Head())

		slog.Info("successor list updated", "list", n.successorList.String())
		return fmt.Errorf("can't retrieve successor %v's predecessor: %v", succ.Identifier(), err)
	}

	// Successor is live
	// Adopt successor list
	err = n.adoptSuccessorList(succ)
	if err != nil {
		return fmt.Errorf("can't adopt successor %v's list %v", succ.Identifier(), err)
	}

	succ, _ = n.Successor()
	if Between(succ_pred.Identifier(), n.Identifier(), succ.Identifier()) {
		n.successorList.SetHead(succ_pred)
		_ = n.adoptSuccessorList(succ_pred)
		n.setSuccessor(succ_pred)
	}

	succ, _ = n.Successor()
	if succ != succStart {
		slog.Info("new successor", "successor", succ)
	}

	return nil
}

// adoptSuccessorList retains the current head of the successor list and copies all but the last entry of p on top
// Not thread safe
func (n *Node) adoptSuccessorList(p node) error {
	if n.Identifier() == p.Identifier() {
		return nil
	}

	newSuccList, err := p.SuccessorList()
	if err != nil {
		return err
	}

	n.successorList.Adopt(newSuccList)

	return nil
}

func (n *Node) SuccessorList() (*SuccessorList, error) {
	return n.successorList, nil
}

// checkPredecessor verifies that the current predecesor is alive. If not, the predecessor is reset.
func (n *Node) checkPredecessor() {
	n.muPred.Lock()
	defer n.muPred.Unlock()
	if n.predecessor != nil && !n.predecessor.Alive() {
		fmt.Println("The predecessor is dead, resetting")
		n.predecessor = nil
	}
}

func (n *Node) Rectify(newPredc node) error {
	n.muPred.Lock()
	defer n.muPred.Unlock()

	pred, _ := n.Predecessor()
	if pred == nil || Between(newPredc.Identifier(), pred.Identifier(), n.Identifier()) {
		slog.Info("accepted rectify", "remote_node", newPredc)
		n.predecessor = newPredc

		n.predecessorGauge.Set(float64(n.predecessor.Identifier()))
		n.operationCount.WithLabelValues("rectify", "success", fmt.Sprint(n.Identifier())).Inc()
	}

	return nil
}

// fixFingers updates the finger table, it is expected to be called repeatedly and updates
// one finger at a time
func (n *Node) fixFingers() {
	if n.nextFinger >= m {
		n.nextFinger = 1
	}

	succ, _, err := n.FindSuccessor(n.id+1<<(n.nextFinger-1), 0)
	if err != nil {
		slog.Warn("failed finger check", "finger", n.nextFinger, "successor", succ)
		n.nextFinger++
		return
	}

	n.muFinger.Lock()
	defer n.muFinger.Unlock()
	n.finger[n.nextFinger] = succ
	n.nextFinger++
}

// FindSuccessor returns the successor node for a given Id by recursively asking the highest
// node in our finger table which comes precedes the given Id
func (n *Node) FindSuccessor(Id Id, pathLength int) (node, int, error) {
	succ, _ := n.Successor()
	if succ == nil {
		n.operationCount.WithLabelValues("findsuccessor", "fail", fmt.Sprint(n.Identifier())).Inc()
		return nil, pathLength, fmt.Errorf("could not find a successor as the node's successor is nil")
	}

	if Between(Id, n.Identifier(), succ.Identifier()+1) {
		return succ, pathLength, nil
	}

	p := n.closestPrecedingNode(Id)
	if p == n {
		return n, pathLength, nil
	}

	n.operationCount.WithLabelValues("findsuccessor", "success", fmt.Sprint(n.Identifier())).Inc()
	return p.FindSuccessor(Id, pathLength+1)
}

// closestPrecedingNode returns the highest entry in the finger table which precedes Id
func (n *Node) closestPrecedingNode(Id Id) node {
	n.muFinger.Lock()
	defer n.muFinger.Unlock()

	for i := m - 1; i >= 0; i-- {
		if n.finger[i] != nil && Between(n.finger[i].Identifier(), n.Identifier(), Id) {
			return n.finger[i]
		}
	}

	return n
}

// String returns a basic string representation of the node for debugging purposes
func (n *Node) String() string {
	var predecessor string = "?"
	var successor string = "?"

	pred, _ := n.Predecessor()
	if pred != nil {
		predecessor = fmt.Sprint(pred.Identifier())
	}

	succ, _ := n.Successor()
	if succ != nil {
		successor = fmt.Sprint(succ.Identifier())
	}

	return fmt.Sprintf("Node(id = %v, predecessor = %v, successor = %v)", n.Identifier(), predecessor, successor)
}

// Alive returns the node's liveness, this is always true for a local node.
func (n *Node) Alive() bool {
	return true
}

func (n *Node) PrometheusRegistry() *prometheus.Registry {
	return n.registry
}
