package main

import (
	"context"
	"fmt"
	"net"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/config"
	"github.com/sicozz/project00/utils"
	"google.golang.org/grpc"
)

type Linker struct {
	proto00.UnimplementedLinkerServer
}

func (s Linker) Info(ctx context.Context, req *proto00.InfoReq) (res *proto00.InfoRes, err error) {
	return &proto00.InfoRes{Version: "0.0", Banner: "PROJECT00", MarketValue: 21}, nil
}

func main() {
	utils.InitLog(utils.DEFAULT_LOG_FILE)
	conf := config.BuildConfig()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	if err != nil {
		utils.Error(fmt.Sprintf("Canot create listener: %s", err))
		return
	}
	server := buildServer()
	err = server.Serve(lis)
	if err != nil {
		utils.Error(fmt.Sprintf("Cannot serve %s", err))
	}
}

func buildServer() *grpc.Server {
	server := grpc.NewServer()
	service := &Linker{}
	proto00.RegisterLinkerServer(server, service)
	return server
}
