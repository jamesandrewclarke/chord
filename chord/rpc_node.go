package chord

import (
	chord_proto "chord_dht/protos/chord"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const TIMEOUT = 10 * time.Second

// RPCNode represents a remote node accessed over the network
type RPCNode struct {
	Address string

	Id Id

	client chord_proto.ChordClient
}

func (n *RPCNode) getConnection() (chord_proto.ChordClient, error) {
	if n.client == nil {
		conn, err := grpc.Dial(n.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("error getting connection: %v\n", err)
			return nil, err
		}

		n.client = chord_proto.NewChordClient(conn)
	}

	return n.client, nil
}

func (n *RPCNode) Identifier() Id {
	return n.Id
}

func (n *RPCNode) Predecessor() (node, error) {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	p, err := chord_client.GetPredecessor(ctx, &chord_proto.PredecessorRequest{})
	if err != nil {
		return nil, err
	}

	newNode := &RPCNode{
		Id:      Id(p.Identifier),
		Address: p.Address,
	}
	SavePeer(newNode)

	return newNode, nil
}

func (n *RPCNode) Successor() (node, error) {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	p, err := chord_client.GetSuccessor(ctx, &chord_proto.SuccessorRequest{})
	if err != nil {
		return nil, err
	}

	newNode := &RPCNode{
		Id:      Id(p.Identifier),
		Address: p.Address,
	}
	SavePeer(newNode)

	return newNode, nil
}

func (n *RPCNode) FindSuccessor(id Id) (node, error) {
	chord_client, err := n.getConnection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	p, err := chord_client.FindSuccessor(ctx, &chord_proto.FindSuccessorRequest{
		Id: int64(id),
	})
	if err != nil {
		return nil, err
	}

	newNode := &RPCNode{
		Id:      Id(p.Identifier),
		Address: p.Address,
	}
	SavePeer(newNode)

	return newNode, nil
}

func (n *RPCNode) Rectify(p node) error {
	chord_client, err := n.getConnection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	node, err := GetPeer(p.Identifier())
	if err != nil {
		return err
	}

	_, err = chord_client.Rectify(ctx, serializePeer(node))

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

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	succListResponse, err := chord_client.SuccessorList(ctx, &chord_proto.SuccessorListRequest{})
	if err != nil {
		return SuccessorList{}, err
	}

	newSuccList := SuccessorList{}

	for i := 0; i < int(succListResponse.NumSuccessors); i++ {
		node := succListResponse.Nodes[i]
		newNode := &RPCNode{
			Address: node.Address,
			Id:      Id(node.Identifier),
		}
		SavePeer(newNode)
		newSuccList.successors[i] = newNode
	}

	return newSuccList, nil
}

func (n *RPCNode) Alive() bool {
	client, _ := n.getConnection()

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	_, err := client.Alive(ctx, &chord_proto.LivenessRequest{})

	if err != nil {
		fmt.Printf("not alive: %v\n", err)
	}
	return err == nil
}

func (n *RPCNode) Announce(port int, addr *string) Id {
	client, _ := n.getConnection()

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	res, err := client.Announce(ctx, &chord_proto.AnnounceRequest{
		Port:    int32(port),
		Address: addr,
	})

	if err != nil {
		return -1
	}

	return Id(res.Identifier)
}

// String returns a basic string representation of the node for debugging purposes
func (n *RPCNode) String() string {
	var addr = "?"
	if n.Address != "" {
		addr = n.Address
	}

	return fmt.Sprintf("RPCNode(id = %v, address = %v)", n.Identifier(), addr)
}
