package chord

import (
	"fmt"
	"log"
	"time"
)

type Id int64

const m = 64

type Node struct {
	id Id

	successor   node
	predecessor node
	finger      [m]node

	nextFinger int
}

type node interface {
	Identifier() Id
	Successor() (node, error)
	Predecessor() (node, error)
	FindSuccessor(Id) (node, error)
	Notify(node) error
}

// CreateNode initialises a single-node Chord ring
func CreateNode(Id Id) *Node {
	n := &Node{
		id:          Id,
		predecessor: nil,
		successor:   nil,
	}

	n.setSuccessor(n)
	n.nextFinger = 1

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
	return n.successor, nil
}

// Start starts the background tasks to stabilize n's pointers and lookup table
func (n *Node) Start() {
	go func() {
		// TODO Configurable intervals for experiments
		for {
			n.stabilize()
			n.fixFingers()
			time.Sleep(2000 * time.Millisecond)
		}
	}()

	go func() {
		for {
			fmt.Println(n)
			time.Sleep(3 * time.Second)
		}
	}()
}

// Join joins a Chord ring containing the node p
func (n *Node) Join(p node) error {
	n.predecessor = nil

	succ, err := p.FindSuccessor(n.Identifier())
	if err != nil {
		return err
	}
	n.setSuccessor(succ)

	return nil
}

// setSuccessor is a safe wrapper method for setting n's immediate successor
func (n *Node) setSuccessor(p node) {
	n.successor = p

	// Finger 0 is also the successor, and should be set every time

	// TODO do we need separate locations
	n.finger[0] = p
}

// Stabilize updates the node's successor and informs them.
// Should be run at a sensible regular interval.
func (n *Node) stabilize() {
	succ, _ := n.Successor()
	succ_pred, _ := succ.Predecessor()

	if succ_pred != nil && between(succ_pred.Identifier(), n.Identifier(), succ.Identifier()) {
		n.setSuccessor(succ_pred)
	}

	err := succ.Notify(n)
	if err != nil {
		log.Printf("error notifying the successor %v", err)
	}
}

// fixFingers updates the finger table, it is expected to be called repeatedly and updates
// one finger at a time
func (n *Node) fixFingers() {
	if n.nextFinger >= m {
		n.nextFinger = 1
	}

	succ, err := n.FindSuccessor(n.id + 1<<(n.nextFinger-1))
	if err != nil {
		log.Printf("error fetching successor for finger %v", n.nextFinger)
		n.nextFinger++
		return
	}

	n.finger[n.nextFinger] = succ
	n.nextFinger++
}

// Notify is called when Node p thinks it is our predecessor
func (n *Node) Notify(p node) error {
	pred, _ := n.Predecessor()
	// If p is between our current predecessor and us, update it
	if pred == nil || between(p.Identifier(), pred.Identifier(), n.Identifier()) {
		n.predecessor = p
	}

	return nil
}

// FindSuccessor returns the successor node for a given Id by recursively asking the highest
// node in our finger table which comes precedes the given Id
func (n *Node) FindSuccessor(Id Id) (node, error) {
	succ, _ := n.Successor()
	if between(Id, n.Identifier(), succ.Identifier()+1) {
		return succ, nil
	}

	p := n.closestPrecedingNode(Id)
	if p == n {
		return n, nil
	}

	return p.FindSuccessor(Id)
}

// closestPrecedingNode returns the highest entry in the finger table which precedes Id
func (n *Node) closestPrecedingNode(Id Id) node {
	for i := m - 1; i >= 0; i-- {
		if n.finger[i] != nil && between(n.finger[i].Identifier(), n.Identifier(), Id) {
			return n.finger[i]
		}
	}

	return n
}

// String returns a basic string representation of the node for debugging purposes
func (n *Node) String() string {
	var predecessor Id = -1

	pred, _ := n.Predecessor()
	if pred != nil {
		predecessor = pred.Identifier()
	}

	succ, _ := n.Successor()
	return fmt.Sprintf("id = %v, predecessor = %v, successor = %v", n.Identifier(), predecessor, succ.Identifier())
}

// For handling circular intervals
func between(id, start, end Id) bool {
	if start < end {
		return id > start && id < end
	}

	return id > start || id < end
}

func IdBetween(id Id, a, b node) bool {
	return between(a.Identifier(), b.Identifier(), id)
}
