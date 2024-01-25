package chord

import (
	chord_proto "chord/protos"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RPCNode represents a remote node accessed over the network
type RPCNode struct {
	Address string

	Id Id
}

func (n *RPCNode) getConnection() (chord_proto.ChordClient, error) {
	conn, err := grpc.Dial(n.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return chord_proto.NewChordClient(conn), err
}

func (n *RPCNode) Identifier() Id {
	return n.Id
}

func (n *RPCNode) Predecessor() (node, error) {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := chord_client.GetPredecessor(ctx, &chord_proto.PredecessorRequest{})
	if err != nil {
		return nil, err
	}

	return &RPCNode{
		Address: p.Address,
	}, nil
}

func (n *RPCNode) Successor() (node, error) {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := chord_client.GetSuccessor(ctx, &chord_proto.SuccessorRequest{})
	if err != nil {
		return nil, err
	}

	return &RPCNode{
		Address: p.Address,
	}, nil
}

func (n *RPCNode) FindSuccessor(Id) (node, error) {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := chord_client.FindSuccessor(ctx, &chord_proto.FindSuccessorRequest{})
	if err != nil {
		return nil, err
	}

	return &RPCNode{
		Address: p.Address,
	}, nil
}

func (n *RPCNode) Notify(node) error {
	chord_client, err := n.getConnection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = chord_client.Notify(ctx, &chord_proto.Node{
		Address:    "127.0.0.1", // get address
		Identifier: idToBytes(n.Identifier()),
	})

	if err != nil {
		return err
	}

	return nil
}

func idToBytes(id Id) []byte {
	return []byte{}
}
