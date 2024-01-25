package chord

import (
	chord_proto "chord/protos"
	"context"
	"log"
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
		Address: p.Address,
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
		Address: p.Address,
	}
}

func (n *RPCNode) FindSuccessor(Id) node {
	chord_client, err := n.getConnection()
	if err != nil {
		log.Println(err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p, err := chord_client.FindSuccessor(ctx, &chord_proto.FindSuccessorRequest{})
	if err != nil {
		log.Println(err)
		return nil
	}

	return &RPCNode{
		Address: p.Address,
	}
}

func (n *RPCNode) Notify(node) {
	chord_client, _ := n.getConnection()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, _ = chord_client.Notify(ctx, &chord_proto.Node{
		Address:    "", // get address
		Identifier: idToBytes(n.Identifier()),
	})

	// TODO error handling
}

func idToBytes(id Id) []byte {
	return []byte{}
}
