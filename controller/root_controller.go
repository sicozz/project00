package controller

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
	"github.com/sicozz/project00/info"
	"github.com/sicozz/project00/logger"
	"github.com/sicozz/project00/network"
	"github.com/sicozz/project00/node"
	"github.com/sicozz/project00/server"
	"google.golang.org/grpc"
)

type RootController struct {
	stm   node.Node
	srv   server.Server
	lis   net.Listener
	gSrv  *grpc.Server
	exitC chan int
	conf  config.Config
}

func NewRootController() (rc RootController, err error) {
	conf := config.BuildConfig()
	logger.InitLog(conf.LogFile)
	addr := conf.GetBindAddr()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot create listener: %v", err))
		return RootController{}, err
	}
	exitC := make(chan int)
	hosts, err := network.DiscoverHosts(conf.HostsFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to discover hosts: %v", err))
		return RootController{}, err
	}
	localhostId, err := network.GetLocalHostIp(hosts)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get own host: %v", err))
		return RootController{}, err
	}
	stm := node.Node(node.NewRaftSTM(hosts, localhostId))
	srv00 := server.NewServer00(stm)
	gSrv := grpc.NewServer()
	proto00.RegisterLinkerServer(gSrv, srv00)
	return RootController{
		stm:   stm,
		srv:   server.Server(srv00),
		lis:   lis,
		gSrv:  gSrv,
		exitC: exitC,
		conf:  conf,
	}, nil
}

func (rc *RootController) Launch() {
	// TODO: Add options for project00
	go rc.handleSignals()
	go rc.startServer()
	go rc.startStateMachine()
	eC := <-rc.exitC
	close(rc.exitC)
	logger.Info(fmt.Sprintf("Exiting %v: %v", info.BANNER, eC))
}

func (rc *RootController) shutDown(exitCode int) {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	rc.srv.Shutdown()
	logger.Info("Shutting down: server stoped")
	rc.gSrv.GracefulStop()
	logger.Info("Shuting down: gRPC server stoped")
	rc.lis.Close()
	logger.Info("Shuting down: listener closed")
	rc.exitC <- exitCode
}

func (rc *RootController) startServer() error {
	logger.Info(fmt.Sprintf("Server listening on port: %v", rc.conf.Port))
	err := rc.gSrv.Serve(rc.lis)
	if err != nil {
		logger.Error(fmt.Sprintf("Server crashed on: %v", err))
		rc.shutDown(1)
	}
	return nil
}

func (rc *RootController) startStateMachine() error {
	logger.Info(fmt.Sprintf("State machine activated"))
	err := rc.stm.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("State machine crashed on: %v", err))
		rc.shutDown(1)
	}
	return nil
}

func (rc *RootController) handleSignals() {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
	rcvdSignal := <-signalC
	close(signalC)
	logger.Info(fmt.Sprintf("Recieved signal: %v", rcvdSignal))
	rc.shutDown(0)
}
