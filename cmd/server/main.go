package main

import (
	"context"
	"fmt"
	"net"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/config"
	"google.golang.org/grpc"
)

type Linker struct {
	proto00.UnimplementedLinkerServer
}

func (s Linker) Info(ctx context.Context, req *proto00.InfoReq) (res *proto00.InfoRes, err error) {
	return &proto00.InfoRes{Version: "0.0", Banner: "PROJECT00", MarketValue: 21}, nil
}

func main() {
	conf := config.BuildConfig()
	// logger
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	if err != nil {
		println("Cannot create listener:", err)
		return
	}
	// server setup
	serverRegistrar := grpc.NewServer()
	service := &Linker{}
	proto00.RegisterLinkerServer(serverRegistrar, service)
	println("Listening on", conf.Port, "...")
	// serve
	err = serverRegistrar.Serve(lis)
	if err != nil {
		println("Cannot create listener %s", err)
	}
}
