package chord

import (
	chord_proto "chord/protos"
	"context"
	"fmt"
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
		fmt.Printf("error getting connection! %v", err)
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
		Id:      Id(p.Identifier),
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
		Id:      Id(p.Identifier),
		Address: p.Address,
	}, nil
}

func (n *RPCNode) FindSuccessor(id Id) (node, error) {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := chord_client.FindSuccessor(ctx, &chord_proto.FindSuccessorRequest{
		Id: int64(id),
	})
	if err != nil {
		return nil, err
	}

	return &RPCNode{
		Id:      Id(p.Identifier),
		Address: p.Address,
	}, nil
}

func (n *RPCNode) Notify(p node) error {
	chord_client, err := n.getConnection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	addr, err := getPeerAddress(p.Identifier())
	if err != nil {
		return err
	}

	_, err = chord_client.Notify(ctx, &chord_proto.Node{
		Address:    addr,
		Identifier: int64(p.Identifier()),
	})

	if err != nil {
		return err
	}

	return nil
}
