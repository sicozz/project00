package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Setup
	utils.InitLog(utils.DEFAULT_LOG_FILE)
	conf := config.BuildConfig()
	addr := conf.GetBindAddr()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		utils.Error(fmt.Sprintf("Cannot create listener: %v", err))
		return
	}
	server := buildServer()
	exitC := make(chan int)

	go handleSignals(server, lis, exitC)
	go launchServer(server, lis, exitC, addr)

	exitCode := <-exitC
	utils.Info(fmt.Sprintf("Exit code %v\n", exitCode))
}

func buildServer() *grpc.Server {
	server := grpc.NewServer()
	service := &Linker{}
	proto00.RegisterLinkerServer(server, service)
	return server
}

func shutdown(srv *grpc.Server, lis net.Listener, exitC chan int, exitCode int) {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	srv.GracefulStop()
	lis.Close()

	utils.Info("Shutdown successfully completed")
	exitC <- exitCode
}

func handleSignals(srv *grpc.Server, lis net.Listener, exitC chan int) {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalC
	utils.Info(fmt.Sprintf("Received signal %v. Initializing graceful shutdown...", sig))
	shutdown(srv, lis, exitC, 0)
}

func launchServer(srv *grpc.Server, lis net.Listener, exitC chan int, addr string) {
	utils.Info(fmt.Sprintf("Serving on: %v...", addr))
	err := srv.Serve(lis)
	if err != nil {
		utils.Error(fmt.Sprintf("Server failed %v", err))
		shutdown(srv, lis, exitC, 1)
	}
}
