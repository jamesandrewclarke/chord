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
	key := in.Key

	if !s.keystore.HasKey(key) {
		chordKey := ChordIdFromString(key)
		successor, err := s.node.FindSuccessor(chordKey)
		if err != nil {
			msg := fmt.Sprintf("Our node does not have this key, and we could not find a node to forward to: %v", err)
			return nil, status.Error(codes.Internal, msg)
		}

		if successor.Identifier() != s.node.Identifier() {
			forwardAddress := stripPort(chord.GetNodeAddress(successor))
			return &dht_proto.GetKeyResponse{
				ForwardNode: &dht_proto.Node{
					Address: forwardAddress,
				},
			}, nil
		}

		return nil, status.Error(codes.Internal, "node does not have this key")
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
	key := in.Key

	// Check if we are actually the successor for this key
	chordKey := ChordIdFromString(in.Key)
	successor, err := s.node.FindSuccessor(chordKey)
	if err != nil {
		msg := fmt.Sprintf("key setting failed, could not verify the node's ownership of the key: %v", err)
		return nil, status.Error(codes.Internal, msg)
	}
	if successor.Identifier() != s.node.Identifier() {
		forwardAddress := stripPort(chord.GetNodeAddress(successor))
		return &dht_proto.SetKeyResponse{
			ForwardNode: &dht_proto.Node{
				Address: forwardAddress,
			},
		}, nil
	}

	err = s.keystore.SetKey(key, in.Value)
	if err != nil {
		return nil, fmt.Errorf("error setting key")
	}

	return &dht_proto.SetKeyResponse{}, nil
}
