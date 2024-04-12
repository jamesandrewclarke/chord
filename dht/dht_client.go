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

func SetKey(address string, key string, value []byte, transfer bool) error {
	fmt.Printf("setting on: %v\n", address)
	client, err := getClient(address)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.SetKey(ctx, &dht_proto.SetKeyRequest{
		Key:      key,
		Value:    value,
		Transfer: transfer,
	})

	if err != nil {
		fmt.Printf("Error setting key: %v\n", err)
		return err
	}

	if !transfer && res.ForwardNode != nil {
		forwardAddr := fmt.Sprintf("%v:%v", res.ForwardNode.Address, DHT_PORT)
		return SetKey(forwardAddr, key, value, transfer)
	}

	return nil
}

func TransferKeys(address string, keys *KeyStore) {
	keys.muKeys.Lock()
	defer keys.muKeys.Unlock()

	for _, v := range keys.Keys {
		v.RLock()
		defer v.RUnlock()

		err := SetKey(address, v.Key, v.Value, true)
		if err != nil {
			fmt.Printf("Error transferring key: %v\n", err)
		}
	}
}
