package chord

import (
	chord_proto "chord/protos"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	local node
	chord_proto.UnimplementedChordServer
}

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
	log.Printf("received predecessor request")

	p, err := s.local.Predecessor()
	if err != nil {
		return &chord_proto.Node{}, err
	}

	addr, err := getPeerAddress(p.Identifier())
	if err != nil {
		return &chord_proto.Node{}, err
	}

	return &chord_proto.Node{
		Address:    addr,
		Identifier: int64(p.Identifier()),
	}, nil
}

func (s *server) GetSuccessor(ctx context.Context, in *chord_proto.SuccessorRequest) (*chord_proto.Node, error) {
	log.Printf("received successor request")

	p, err := s.local.Successor()
	if err != nil {
		return nil, err
	}

	addr, err := getPeerAddress(p.Identifier())
	if err != nil {
		return nil, err
	}

	return &chord_proto.Node{
		Address:    addr,
		Identifier: int64(p.Identifier()),
	}, nil
}

func (s *server) FindSuccessor(ctx context.Context, in *chord_proto.FindSuccessorRequest) (*chord_proto.Node, error) {
	log.Printf("received findsuccessor request")

	lookupID := in.Id
	p, err := s.local.FindSuccessor(Id(lookupID))
	if err != nil {
		return nil, err
	}

	foundID := p.Identifier()
	addr, err := getPeerAddress(foundID)
	if err != nil {
		fmt.Printf("couldn't find peer address %v", err)
		return nil, err
	}

	return &chord_proto.Node{
		Identifier: int64(foundID),
		Address:    addr,
	}, nil
}

func (s *server) Rectify(ctx context.Context, in *chord_proto.Node) (*chord_proto.RectifyResponse, error) {
	log.Printf("received rectify request from %v", in.Identifier)

	node := &RPCNode{
		Address: in.Address,
		Id:      Id(in.Identifier),
	}

	s.local.Rectify(node)

	return &chord_proto.RectifyResponse{}, nil
}

func (s *server) SuccessorList(ctx context.Context, in *chord_proto.SuccessorListRequest) (*chord_proto.SuccessorListResponse, error) {
	succ_list, _ := s.local.SuccessorList()
	response := &chord_proto.SuccessorListResponse{}
	response.Nodes = make([]*chord_proto.Node, r)
	for i, succ := range succ_list.successors {
		if succ == nil {
			break
		}
		addr, _ := getPeerAddress(succ.Identifier())
		response.NumSuccessors++
		response.Nodes[i] = &chord_proto.Node{
			Address:    addr,
			Identifier: int64(succ.Identifier()),
		}
	}

	return response, nil
}
