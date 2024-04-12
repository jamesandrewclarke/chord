package dht

import (
	dht_proto "chord_dht/protos/dht"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DHT_PORT = 8081

func getClient(address string) (dht_proto.DHTClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("error getting connection: %v\n", err)
		return nil, err
	}

	client := dht_proto.NewDHTClient(conn)
	return client, nil
}

func SetKey(address string, key string, value []byte) error {
	fmt.Printf("setting on: %v\n", address)
	client, err := getClient(address)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := client.SetKey(ctx, &dht_proto.SetKeyRequest{
		Key:   key,
		Value: value,
	})

	if err != nil {
		fmt.Printf("Error setting key: %v\n", err)
		return err
	}

	if res.ForwardNode != nil {
		forwardAddr := fmt.Sprintf("%v:%v", res.ForwardNode.Address, DHT_PORT)
		return SetKey(forwardAddr, key, value)
	}

	return nil
}
