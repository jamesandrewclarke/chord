package chord

import (
	chord_proto "chord/protos"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// RPCNode represents a remote node accessed over the network
type RPCNode struct {
	ipv4 string

	id Id
}

func (n *RPCNode) getConnection() (chord_proto.ChordClient, error) {
	conn, err := grpc.Dial(n.ipv4)
	if err != nil {
		return nil, err
	}

	return chord_proto.NewChordClient(conn), err
}

func (n *RPCNode) Identifier() Id {
	return n.id
}

func (n *RPCNode) Predecessor() node {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := chord_client.GetPredecessor(ctx, &chord_proto.PredecessorRequest{})
	if err != nil {
		return nil
	}

	return &RPCNode{
		ipv4: fmt.Sprintf("%v", p.GetIpv4()),
	}
}

func (n *RPCNode) Successor() node {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := chord_client.GetSuccessor(ctx, &chord_proto.SuccessorRequest{})
	if err != nil {
		return nil
	}

	return &RPCNode{
		ipv4: fmt.Sprintf("%v", p.GetIpv4()),
	}
}

func (n *RPCNode) FindSuccessor(Id) node {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := chord_client.FindSuccessor(ctx, &chord_proto.FindSuccessorRequest{})
	if err != nil {
		return nil
	}

	return &RPCNode{
		ipv4: fmt.Sprintf("%v", p.GetIpv4()),
	}
}

func (n *RPCNode) Notify(node) {
}
