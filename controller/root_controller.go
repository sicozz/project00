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
	"github.com/sicozz/project00/server"
	"github.com/sicozz/project00/statemachine"
	"github.com/sicozz/project00/utils"
	"google.golang.org/grpc"
)

type RootController struct {
	conf     config.Config
	exitC    chan int
	srv      *grpc.Server
	server00 *server.Server00
	stm00    *statemachine.StateMachine00
	lis      net.Listener
}

func NewRootController() (rc RootController, err error) {
	conf := config.BuildConfig()
	addr := conf.GetBindAddr()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		utils.Error(fmt.Sprintf("Cannot create listener: %v", err))
		return RootController{}, err
	}
	srv := grpc.NewServer()
	server00 := server.NewServer00()
	proto00.RegisterLinkerServer(srv, server00)
	stm00 := statemachine.NewStateMachine00(server00.STMCh)
	exitC := make(chan int)
	return RootController{
		conf:     conf,
		exitC:    exitC,
		srv:      srv,
		server00: server00,
		stm00:    stm00,
		lis:      lis,
	}, nil
}

func (rc *RootController) Launch() {
	// TODO: Add options for project00
	go rc.handleSignals()
	go rc.startServer()
	go rc.stm00.Run()
	eC := <-rc.exitC
	close(rc.exitC)
	utils.Info(fmt.Sprintf("Exiting project00: %v", eC))
	rc.server00.Shutdown()
}

func (rc *RootController) shutDown(exitCode int) {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	rc.srv.GracefulStop()
	utils.Info("Shuting down: Server stoped")
	rc.lis.Close()
	utils.Info("Shuting down: Listener closed")
	rc.exitC <- exitCode
}

func (rc *RootController) startServer() error {
	utils.Info(fmt.Sprintf("Server listening on port: %v", rc.conf.Port))
	err := rc.srv.Serve(rc.lis)
	if err != nil {
		utils.Error(fmt.Sprintf("Server crashed on: %v", err))
		rc.shutDown(1)
	}
	return nil
}

func (rc *RootController) handleSignals() {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
	rcvdSignal := <-signalC
	close(signalC)
	utils.Info(fmt.Sprintf("Recieved signal: %v", rcvdSignal))
	rc.shutDown(0)
}
