package server

import (
	"context"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/handler"
)

type LinkerService struct {
	proto00.UnimplementedLinkerServer
}

func (s *LinkerService) Info(ctx context.Context, req *proto00.InfoReq) (res *proto00.InfoRes, err error) {
	return handler.QueryGetServiceInfo()
}

func (s *LinkerService) Subscribe(req *proto00.SubscribeReq, stream proto00.Linker_SubscribeServer) error {
	return handler.HandleSubscription(stream)
}
