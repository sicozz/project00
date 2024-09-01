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
	stm   statemachine.StateMachine
	srv   server.Server
	lis   net.Listener
	gSrv  *grpc.Server
	exitC chan int
	conf  config.Config
}

func NewRootController() (rc RootController, err error) {
	conf := config.BuildConfig()
	addr := conf.GetBindAddr()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		utils.Error(fmt.Sprintf("Cannot create listener: %v", err))
		return RootController{}, err
	}
	exitC := make(chan int)
	stmAndSrvCh := make(chan string)
	stm00 := statemachine.NewStateMachine00(stmAndSrvCh)
	srv00 := server.NewServer00(stmAndSrvCh)
	gSrv := grpc.NewServer()
	proto00.RegisterLinkerServer(gSrv, srv00)
	return RootController{
		stm:   statemachine.StateMachine(stm00),
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
	utils.Info(fmt.Sprintf("Exiting %v: %v", utils.BANNER, eC))
}

func (rc *RootController) shutDown(exitCode int) {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	rc.srv.Shutdown()
	utils.Info("Shutting down: server stoped")
	rc.gSrv.GracefulStop()
	utils.Info("Shuting down: gRPC server stoped")
	rc.lis.Close()
	utils.Info("Shuting down: listener closed")
	rc.exitC <- exitCode
}

func (rc *RootController) startServer() error {
	utils.Info(fmt.Sprintf("Server listening on port: %v", rc.conf.Port))
	err := rc.gSrv.Serve(rc.lis)
	if err != nil {
		utils.Error(fmt.Sprintf("Server crashed on: %v", err))
		rc.shutDown(1)
	}
	return nil
}

func (rc *RootController) startStateMachine() error {
	utils.Info(fmt.Sprintf("State machine activated"))
	err := rc.stm.Run()
	if err != nil {
		utils.Error(fmt.Sprintf("State machine crashed on: %v", err))
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
