package chord

import (
	"fmt"
	"time"
)

type Id int

const m = 64

type Node struct {
	Id          Id
	Successor   *Node
	Predecessor *Node
	Finger      [m]*Node

	nextFinger int
}

// TODO use a generic instead of 'int' so we can change it later for a different type
type node interface {
	Id() int
	Successor() int
	Predecessor() int
	FindSuccessor(int) node
	Notify(node)
}

// CreateNode initialises a single-node Chord ring
func CreateNode(Id Id) *Node {
	n := &Node{
		Id:          Id,
		Predecessor: nil,
		Successor:   nil,
	}

	n.Successor = n
	n.nextFinger = 1

	return n
}

func (n *Node) Start() {
	go func() {
		for {
			n.stabilize()
			// n.fixFingers()
			fmt.Println(n)
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

// Join joins a Chord ring containing the node p
func (n *Node) Join(p *Node) {
	n.Predecessor = nil
	n.Successor = p.FindSuccessor(n.Id)
}

// Stabilize updates the node's successor and informs them.
// Should be run at a sensible regular interval.
func (n *Node) stabilize() {
	x := n.Successor.Predecessor
	if x != nil && between(x.Id, n.Id, n.Successor.Id) {
		n.Successor = x
	}

	n.Successor.Notify(n)
}

func (n *Node) fixFingers() {
	if n.nextFinger >= m {
		n.nextFinger = 1
	}

	n.Finger[n.nextFinger] = n.FindSuccessor(n.Id + 1<<(n.nextFinger-1))
	n.nextFinger++
}

// Notify is called when Node p thinks it is our predecessor
func (n *Node) Notify(p *Node) {
	// If p is between our current predecessor and us, update it
	if n.Predecessor == nil || between(p.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = p
	}
}

// FindSuccessorIterative returns the node succeeding a given ID without the finger table
func (n *Node) FindSuccessorIterative(Id Id) *Node {
	if n == n.Successor {
		return n
	}

	if betweenIncStart(Id, n.Id, n.Successor.Id) {
		return n.Successor
	} else {
		// Just forward the query around the circle until we find it
		return n.Successor.FindSuccessorIterative(Id)
	}
}

func (n *Node) FindSuccessor(Id Id) *Node {
	if between(Id, n.Id, n.Successor.Id) {
		return n.Successor
	}

	p := n.ClosestPrecedingNode(Id)
	if p == n {
		return n
	}

	return p.FindSuccessor(Id)
}

func (n *Node) ClosestPrecedingNode(Id Id) *Node {
	for i := m - 1; i >= 0; i-- {
		if n.Finger[i] != nil && between(n.Finger[i].Id, n.Id, Id) {
			return n.Finger[i]
		}
	}

	return n
}

func (n *Node) String() string {
	var predecessor Id = -1
	if n.Predecessor != nil {
		predecessor = n.Predecessor.Id
	}
	var successor Id = -1
	successor = n.Successor.Id
	return fmt.Sprintf("id = %v, predecessor = %v, successor = %v", n.Id, predecessor, successor)
}

// For handling circular intervals
func between(id, start, end Id) bool {
	if start < end {
		return id > start && id < end
	}

	return id > start || id < end
}

func betweenIncStart(id, start, end Id) bool {
	if start < end {
		return id >= start && id < end
	}

	return id >= start || id < end
}
