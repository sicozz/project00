package main

import (
	"context"
	"net"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"google.golang.org/grpc"
)

type Linker struct {
	proto00.UnimplementedLinkerServer
}

func (s Linker) Info(ctx context.Context, req *proto00.InfoReq) (res *proto00.InfoRes, err error) {
	return &proto00.InfoRes{Version: "0.0", Banner: "PROJECT00", MarketValue: 21}, nil
}

func main() {
	lis, err := net.Listen("tcp", "[::]:50051")
	if err != nil {
		println("Cannot create listener %s", err)
		return
	}
	serverRegistrar := grpc.NewServer()
	service := &Linker{}
	proto00.RegisterLinkerServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		println("Cannot create listener %s", err)
	}
}
