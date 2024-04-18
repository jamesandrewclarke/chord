package dht

import (
	"chord_dht/chord"
	dht_proto "chord_dht/protos/dht"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	node *chord.Node

	shutdown chan struct{}
	wg       *sync.WaitGroup

	keystore *KeyStore
	dht_proto.UnimplementedDHTServer
}

func StartDHT(node *chord.Node, port int) *Server {
	s := grpc.NewServer()

	dht := &Server{
		node:     node,
		keystore: CreateKeyStore(node.Identifier()),
		shutdown: make(chan struct{}),
		wg:       new(sync.WaitGroup),
	}

	prometheus.MustRegister(dht.keystore.Registry)
	dht_proto.RegisterDHTServer(s, dht)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		panic(err)
	}

	dht.wg.Add(1)
	go func() {
		defer dht.wg.Done()

		keyCheckTicker := time.NewTicker(3 * time.Second)
		defer keyCheckTicker.Stop()

		for {
			select {
			case <-keyCheckTicker.C:
				dht.CheckKeys()

			case <-dht.shutdown:
				fmt.Println("Stopping...")
				node.Stop()

				succ, _ := node.Successor()
				if succ != nil && succ != node {
					fmt.Printf("Transferring keys to %v\n", succ)
					addr := fmt.Sprintf("%v:%v", stripPort(chord.GetNodeAddress(succ)), DHT_PORT)
					TransferKeys(addr, dht.keystore)
				}
				return
			}
		}
	}()

	go func() {
		log.Printf("DHT server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return dht
}

func (s *Server) Stop() {
	fmt.Println("Starting DHT graceful shutdown...")
	close(s.shutdown)
	s.wg.Wait()
	fmt.Println("Done")
}

func (s *Server) CheckKeys() {
	for _, v := range s.keystore.Keys {
		v.RLock()

		// First check if the key is between the current predecessor
		// and us, if it is, then continue
		pred, err := s.node.Predecessor()
		if err != nil {
			fmt.Println("key check failed, no predecessor")
			v.RUnlock()
			continue
		}

		if chord.Between(v.Id, pred.Identifier()+1, s.node.Identifier()) {
			v.RUnlock()
			continue
		}

		fmt.Printf("Transferring key: %v\n", v.Id)
		predAddr := fmt.Sprintf("%v:%v", stripPort(chord.GetNodeAddress(pred)), DHT_PORT)
		err = SetKey(predAddr, v.Key, v.Value, true)
		v.RUnlock()

		if err != nil {
			fmt.Printf("error transferring key... %v", err)
		} else {
			fmt.Println("successfully transferred key, deleting...")
			s.keystore.DeleteKey(v.Key)
			fmt.Println("deleted")
		}
	}
}

func (s *Server) GetKey(ctx context.Context, in *dht_proto.GetKeyRequest) (*dht_proto.GetKeyResponse, error) {
	key := in.Key

	fmt.Printf("Received GetKey for %v\n", key)

	if !s.keystore.HasKey(key) {
		chordKey := ChordIdFromString(key)
		successor, pathLength, err := s.node.FindSuccessor(chordKey, 0)
		fmt.Printf("Path length: %v\n", pathLength)
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
				PathLength: int32(pathLength),
			}, nil
		}

		return nil, status.Error(codes.Internal, "node does not have this key")
	}

	value, err := s.keystore.GetKey(key)
	if err != nil {
		return nil, fmt.Errorf("error whilst retrieving key")
	}

	return &dht_proto.GetKeyResponse{
		Value:      value,
		PathLength: 0,
	}, nil
}

func (s *Server) SetKey(ctx context.Context, in *dht_proto.SetKeyRequest) (*dht_proto.SetKeyResponse, error) {
	key := in.Key

	fmt.Printf("Received SetKey for %v\n", key)
	// Check if we are actually the successor for this key
	chordKey := ChordIdFromString(in.Key)
	if !in.Transfer {
		successor, _, err := s.node.FindSuccessor(chordKey, 0)

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
	}

	err := s.keystore.SetKey(key, in.Value)
	if err != nil {
		return nil, fmt.Errorf("error setting key")
	}

	return &dht_proto.SetKeyResponse{}, nil
}
