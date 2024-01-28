package chord

import (
	"fmt"
	"log"
	"time"
)

type Id int

const m = 64

type Node struct {
	id Id

	successor   node
	predecessor node
	finger      [m]node

	nextFinger int
}

// TODO use a generic instead of 'int' so we can change it later for a different type
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

	n.successor = n
	n.nextFinger = 1

	return n
}

func (n *Node) Identifier() Id {
	return n.id
}

func (n *Node) Predecessor() (node, error) {
	return n.predecessor, nil
}

func (n *Node) Successor() (node, error) {
	return n.successor, nil
}

func (n *Node) Start() {
	go func() {
		for {
			n.stabilize()
			n.fixFingers()
			time.Sleep(250 * time.Millisecond)
		}
	}()

	go func() {
		for {
			fmt.Println(n)
			time.Sleep(1 * time.Second)
		}
	}()
}

// Join joins a Chord ring containing the node p
func (n *Node) Join(p node) {
	n.predecessor = nil

	succ, _ := p.FindSuccessor(n.Identifier())
	n.successor = succ
}

// Stabilize updates the node's successor and informs them.
// Should be run at a sensible regular interval.
func (n *Node) stabilize() {
	succ, _ := n.Successor()
	succ_pred, err := succ.Predecessor()
	if err != nil {
		log.Printf("error during stabilization, %v", err)
		return
	}

	if succ_pred != nil && between(succ_pred.Identifier(), n.Identifier(), succ.Identifier()) {
		n.successor = succ_pred
	}

	err = succ.Notify(n)
	if err != nil {
		log.Printf("error notifying the successor %v", err)
	}
}

func (n *Node) fixFingers() {
	if n.nextFinger >= m {
		n.nextFinger = 1
	}

	succ, err := n.FindSuccessor(n.id + 1<<(n.nextFinger-1)) // TODO fix this to wrap around to the start of the circle
	if err != nil {
		log.Printf("error fetching successor for finger %v", n.nextFinger)
		return
	}

	n.finger[n.nextFinger] = succ
	// should do tests to verify this
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

func (n *Node) closestPrecedingNode(Id Id) node {
	for i := m - 1; i >= 0; i-- {
		if n.finger[i] != nil && between(n.finger[i].Identifier(), n.Identifier(), Id) {
			return n.finger[i]
		}
	}

	return n
}

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
