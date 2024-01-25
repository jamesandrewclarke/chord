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

	return &chord_proto.Node{}, nil
}

func (s *server) GetSuccessor(ctx context.Context, in *chord_proto.SuccessorRequest) (*chord_proto.Node, error) {
	log.Printf("received successor request")

	return &chord_proto.Node{}, nil
}

func (s *server) FindSuccessor(ctx context.Context, in *chord_proto.FindSuccessorRequest) (*chord_proto.Node, error) {
	log.Printf("received findsuccessor request")

	return &chord_proto.Node{}, nil
}

func (s *server) Notify(ctx context.Context, in *chord_proto.Node) (*chord_proto.NotifyResponse, error) {
	log.Printf("received notify request")

	// stub
	// TODO actually update proto to be able to map an address to a node

	p, _ := peer.FromContext(ctx)

	node := &RPCNode{
		Address: p.Addr.String(),
		Id:      0,
	}

	s.local.Notify(node)

	return &chord_proto.NotifyResponse{}, nil
}
