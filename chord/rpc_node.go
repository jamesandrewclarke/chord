package chord

// RPCNode represents a remote node accessed over the network
type RPCNode struct {
	// ip address

	id Id
}

func (n *RPCNode) Identifier() Id {
	return n.id
}

func (n *RPCNode) Predecessor() node {
	return nil
}

func (n *RPCNode) Successor() node {
	return nil
}

func (n *RPCNode) FindSuccessor(Id) node {
	return nil
}

func (n *RPCNode) Notify(node) {
}
