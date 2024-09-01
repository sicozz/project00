package main

import (
	"context"
	"io"
	"log"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := proto00.NewLinkerClient(conn)
	req := &proto00.SubscribeReq{}
	stream, err := client.Subscribe(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Subscribe RPC: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving heartbeat: %v", err)
		}
		log.Printf("Heartbeat [%v]", resp)
	}
}
