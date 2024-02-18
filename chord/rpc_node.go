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

func (n *RPCNode) Rectify(p node) error {
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

	_, err = chord_client.Rectify(ctx, &chord_proto.Node{
		Address:    addr,
		Identifier: int64(p.Identifier()),
	})

	if err != nil {
		return err
	}

	return nil
}

func (n *RPCNode) SuccessorList() (SuccessorList, error) {
	chord_client, err := n.getConnection()
	if err != nil {
		return SuccessorList{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	succListResponse, err := chord_client.SuccessorList(ctx, &chord_proto.SuccessorListRequest{})
	if err != nil {
		return SuccessorList{}, err
	}

	newSuccList := SuccessorList{}

	for i := 0; i < int(succListResponse.NumSuccessors); i++ {
		node := succListResponse.Nodes[i]
		addr, _ := getPeerAddress(Id(node.Identifier))
		newSuccList.successors[i] = &RPCNode{
			Address: addr,
			Id:      Id(node.Identifier),
		}
	}

	return newSuccList, nil
}

func (n *RPCNode) Alive() bool {
	client, _ := n.getConnection()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.GetSuccessor(ctx, &chord_proto.SuccessorRequest{})
	if err != nil {
		return false
	}

	return true
}

// String returns a basic string representation of the node for debugging purposes
func (n *RPCNode) String() string {
	var predecessor Id = -1

	pred, _ := n.Predecessor()
	if pred != nil {
		predecessor = pred.Identifier()
	}

	succ, _ := n.Successor()
	return fmt.Sprintf("id = %v, predecessor = %v, successor = %v", n.Identifier(), predecessor, succ.Identifier())
}
