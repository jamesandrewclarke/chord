package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "chord/protos"
)

var (
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	key   = flag.String("key", "testing", "the key to set")
	value = flag.String("value", "hellohello", "the value to set")
)

func main() {
	fmt.Println("Hello, World!")

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewKeyValueClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SetValue(ctx, &pb.KeyPair{
		Key:   *key,
		Value: *value,
	})
	if err != nil {
		log.Fatalf("could not set key: %v", err)
	}

	log.Printf("set key %v: %v", *key, *value)

	r, err = c.GetValue(ctx, &pb.KeyPair{
		Key:   *key,
		Value: *value,
	})
	if err != nil {
		log.Fatalf("could not set key: %v", err)
	}

	log.Printf("%v: %v", *key, r.GetValue())
}
