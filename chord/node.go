package chord

import (
	"fmt"
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
	Successor() node
	Predecessor() node
	FindSuccessor(Id) node
	Notify(node)
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

func (n *Node) Predecessor() node {
	return n.predecessor
}

func (n *Node) Successor() node {
	return n.successor
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
func (n *Node) Join(p *Node) {
	n.predecessor = nil
	n.successor = p.FindSuccessor(n.Identifier())
}

// Stabilize updates the node's successor and informs them.
// Should be run at a sensible regular interval.
func (n *Node) stabilize() {
	x := n.Successor().Predecessor()
	if x != nil && between(x.Identifier(), n.Identifier(), n.Successor().Identifier()) {
		n.successor = x
	}

	n.Successor().Notify(n)
}

func (n *Node) fixFingers() {
	if n.nextFinger >= m {
		n.nextFinger = 1
	}

	n.finger[n.nextFinger] = n.FindSuccessor(n.id + 1<<(n.nextFinger-1)) // TODO fix this to wrap around to the start of the circle
	// should do tests to verify this
	n.nextFinger++
}

// Notify is called when Node p thinks it is our predecessor
func (n *Node) Notify(p node) {
	// If p is between our current predecessor and us, update it
	if n.Predecessor() == nil || between(p.Identifier(), n.Predecessor().Identifier(), n.Identifier()) {
		n.predecessor = p
	}
}

func (n *Node) FindSuccessor(Id Id) node {
	if between(Id, n.Identifier(), n.Successor().Identifier()) {
		return n.Successor()
	}

	p := n.closestPrecedingNode(Id)
	if p == n {
		return n
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
	if n.Predecessor() != nil {
		predecessor = n.Predecessor().Identifier()
	}

	successor := n.Successor().Identifier()
	return fmt.Sprintf("id = %v, predecessor = %v, successor = %v", n.Identifier(), predecessor, successor)
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
