package chord

import (
	chord_proto "chord/protos"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type server struct {
	local node
	chord_proto.UnimplementedChordServer
}

var externalAddress string

func StartServer(node node, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	chord_proto.RegisterChordServer(s, &server{local: node})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) GetPredecessor(ctx context.Context, in *chord_proto.PredecessorRequest) (*chord_proto.Node, error) {
	p, err := s.local.Predecessor()
	if err != nil {
		fmt.Printf("%v\n", err)
		return &chord_proto.Node{}, err
	}

	node, err := GetPeer(p.Identifier())
	if err != nil {
		return &chord_proto.Node{}, err
	}

	return serializePeer(node), nil
}

func (s *server) GetSuccessor(ctx context.Context, in *chord_proto.SuccessorRequest) (*chord_proto.Node, error) {
	p, err := s.local.Successor()
	if err != nil {
		return nil, err
	}

	node, err := GetPeer(p.Identifier())
	if err != nil {
		return nil, err
	}

	return serializePeer(node), nil
}

func (s *server) FindSuccessor(ctx context.Context, in *chord_proto.FindSuccessorRequest) (*chord_proto.Node, error) {
	lookupID := in.Id
	p, err := s.local.FindSuccessor(Id(lookupID))
	if err != nil {
		return nil, err
	}

	foundID := p.Identifier()
	node, err := GetPeer(foundID)
	if err != nil {
		fmt.Printf("couldn't find peer %v", err)
		return nil, err
	}

	return serializePeer(node), nil
}

func (s *server) Rectify(ctx context.Context, in *chord_proto.Node) (*chord_proto.RectifyResponse, error) {
	node := &RPCNode{
		Address: in.Address,
		Id:      Id(in.Identifier),
	}

	// Definitely validate
	SavePeer(node)
	s.local.Rectify(node)

	return &chord_proto.RectifyResponse{}, nil
}

func (s *server) SuccessorList(ctx context.Context, in *chord_proto.SuccessorListRequest) (*chord_proto.SuccessorListResponse, error) {
	succ_list, _ := s.local.SuccessorList()
	response := &chord_proto.SuccessorListResponse{}
	response.Nodes = make([]*chord_proto.Node, SUCCESSOR_LIST_SIZE)
	for i, succ := range succ_list.successors {
		if succ == nil {
			break
		}

		response.NumSuccessors++
		response.Nodes[i] = serializePeer(succ)
	}

	return response, nil
}

func (s *server) Announce(ctx context.Context, in *chord_proto.AnnounceRequest) (*chord_proto.Node, error) {
	// Take the announcement message and update the directory

	// Extract the IP address from the call
	// Add to directory along with user supplied port

	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve peer")
	}

	var host string
	if in.Address == nil {
		// User has not supplied optional return address,
		// so retrieve it from the network layer
		host, _, _ = net.SplitHostPort(p.Addr.String())
	} else {
		host = *in.Address
	}

	endpointAddress := fmt.Sprintf("[%s]:%d", host, in.Port)
	id := IdentifierFromAddress(endpointAddress)

	// Update the directory with the new peer
	// perhaps some error handling
	log.Printf("Discovered new peer %v: %v", id, endpointAddress)

	newNode := &RPCNode{
		Id:      id,
		Address: endpointAddress,
	}
	SavePeer(newNode)

	return &chord_proto.Node{
		Address:    endpointAddress,
		Identifier: int64(id),
	}, nil
}

func (s *server) Alive(ctx context.Context, in *chord_proto.LivenessRequest) (*chord_proto.LivenessResponse, error) {
	return &chord_proto.LivenessResponse{}, nil
}

func SetExternalAddress(addr string) {
	externalAddress = addr
}

func serializePeer(node node) *chord_proto.Node {
	res := chord_proto.Node{
		Identifier: int64(node.Identifier()),
	}

	switch v := node.(type) {
	case *Node:
		res.Address = externalAddress

	case *RPCNode:
		res.Address = v.Address
	}

	return &res
}
