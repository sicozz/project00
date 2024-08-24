package server

import (
	"context"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/handler"
)

type LinkerService struct {
	proto00.UnimplementedLinkerServer
}

func (s LinkerService) Info(ctx context.Context, req *proto00.InfoReq) (res *proto00.InfoRes, err error) {
	programInfo := handler.QueryGetServiceInfo()
	return &proto00.InfoRes{
		Version: programInfo.Version,
		Banner:  programInfo.Banner,
	}, nil
}
