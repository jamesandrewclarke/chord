package main

import (
	pb "chord/protos"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
	keys map[string]string
)

type server struct {
	pb.UnimplementedKeyValueServer
}

func (s *server) SetValue(ctx context.Context, in *pb.KeyPair) (*pb.KeyPair, error) {
	log.Printf("Received: %v, %v", in.GetKey(), in.GetValue())

	keys[in.Key] = in.Value
	return in, nil
}

func (s *server) GetValue(ctx context.Context, in *pb.KeyPair) (*pb.KeyPair, error) {
	log.Printf("Received request to get: %v", in.GetKey())

	key := in.GetKey()
	if val, ok := keys[key]; ok {
		in.Value = val
		return in, nil
	}

	return &pb.KeyPair{Key: "", Value: ""}, nil
}

func init() {
	keys = make(map[string]string)
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterKeyValueServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
