package dht

import (
	"chord_dht/chord"
	dht_proto "chord_dht/protos/dht"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	node *chord.Node

	keystore keystore
	dht_proto.UnimplementedDHTServer
}

func StartDHT(node *chord.Node, port int) {
	s := grpc.NewServer()
	dht_proto.RegisterDHTServer(s, &server{
		node:     node,
		keystore: CreateKeyStore(),
	})

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		panic(err)
	}

	log.Printf("DHT server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}

}

func (s *server) GetKey(ctx context.Context, in *dht_proto.GetKeyRequest) (*dht_proto.GetKeyResponse, error) {
	key := chord.Id(in.Key)
	if !s.keystore.HasKey(key) {
		return nil, fmt.Errorf("key not found in this node")
	}

	value, err := s.keystore.GetKey(key)
	if err != nil {
		return nil, fmt.Errorf("error whilst retrieving key")
	}

	return &dht_proto.GetKeyResponse{
		Value: value,
	}, nil
}

func (s *server) SetKey(ctx context.Context, in *dht_proto.SetKeyRequest) (*dht_proto.SetKeyResponse, error) {
	key := chord.Id(in.Key)

	// Calculate the key ourselves as an integrity check
	hash := chord.Hash(in.Value)
	keyRecalculated := chord.IdentifierFromBytes(hash)
	if keyRecalculated != key {
		msg := fmt.Sprintf("integrity check failed, provided key: %v, actual key: %v", key, keyRecalculated)
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	// Check if we are actually the successor for this key
	successor, err := s.node.FindSuccessor(key)
	if err != nil {
		msg := fmt.Sprintf("key setting failed, could not verify the node's ownership of the key: %v", err)
		return nil, status.Error(codes.Internal, msg)
	}
	if successor.Identifier() != s.node.Identifier() {
		msg := fmt.Sprintf("rejected key, node %v is the successor of the provided key", successor.Identifier())
		return nil, status.Error(codes.Canceled, msg)
	}

	err = s.keystore.SetKey(key, in.Value)
	if err != nil {
		return nil, fmt.Errorf("error setting key")
	}

	return &dht_proto.SetKeyResponse{}, nil
}
