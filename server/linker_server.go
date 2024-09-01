package server

import (
	"context"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/handler"
	"google.golang.org/grpc"
)

type Server00 struct {
	STMCh chan string
	*grpc.Server
	proto00.UnimplementedLinkerServer
}

func NewServer00() *Server00 {
	stmCh := make(chan string)
	return &Server00{STMCh: stmCh}
}

func (s *Server00) Info(ctx context.Context, req *proto00.InfoReq) (res *proto00.InfoRes, err error) {
	return handler.QueryGetServiceInfo()
}

func (s *Server00) Subscribe(req *proto00.SubscribeReq, stream proto00.Linker_SubscribeServer) error {
	return handler.HandleSubscription(stream, s.STMCh)
}

func (s *Server00) Shutdown() {
	close(s.STMCh)
}
