package server

import (
	"context"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/handler"
)

type Server interface {
	Shutdown()
}

type Server00 struct {
	stmC chan string
	proto00.UnimplementedLinkerServer
}

func NewServer00(stmC chan string) *Server00 {
	return &Server00{stmC: stmC}
}

func (s *Server00) Info(ctx context.Context, req *proto00.InfoReq) (res *proto00.InfoRes, err error) {
	return handler.QueryGetServiceInfo()
}

func (s *Server00) Subscribe(req *proto00.SubscribeReq, stream proto00.Linker_SubscribeServer) error {
	return handler.HandleSubscription(stream, s.stmC)
}

func (s *Server00) Shutdown() {
	close(s.stmC)
}
