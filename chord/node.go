package chord

import "fmt"

type Node struct {
	Id          int
	Successor   *Node
	Predecessor *Node
}

// CreateNode initialises a single-node Chord ring
func CreateNode(ID int) *Node {
	n := &Node{
		Id:          ID,
		Predecessor: nil,
		Successor:   nil,
	}

	n.Successor = n

	return n
}

// Join joins a Chord ring containing the node p
func (n *Node) Join(p *Node) {
	n.Predecessor = nil
	n.Successor = p.FindSuccessor(n.Id)
}

func (n *Node) Stabilize() {
	x := n.Successor.Predecessor
	if x != nil && between(x.Id, n.Id, n.Successor.Id) {
		fmt.Println("updating")
		n.Successor = x
	}

	n.Successor.Notify(n)
}

// Notify is called when Node p thinks it is our predecessor
func (n *Node) Notify(p *Node) {
	// If p is between our current predecessor and us, update it
	if n.Predecessor == nil || between(p.Id, n.Predecessor.Id, n.Id) {
		n.Predecessor = p
	}
}

// FindSuccessor returns the node succeeding a given ID
func (n *Node) FindSuccessor(Id int) *Node {
	if n == n.Successor {
		return n
	}

	if betweenIncStart(Id, n.Id, n.Successor.Id) {
		return n.Successor
	} else {
		// Just forward the query around the circle until we find it
		return n.Successor.FindSuccessor(Id)
	}
}

func (n *Node) String() string {
	predecessor := -1
	if n.Predecessor != nil {
		predecessor = n.Predecessor.Id
	}
	successor := -1
	successor = n.Successor.Id
	return fmt.Sprintf("id = %v, predecessor = %v, successor = %v", n.Id, predecessor, successor)
}

// For handling circular intervals
func between(id, start, end int) bool {
	if start < end {
		return id > start && id < end
	}

	return id > start || id < end
}

func betweenIncStart(id, start, end int) bool {
	if start < end {
		return id >= start && id < end
	}

	return id >= start || id < end
}
