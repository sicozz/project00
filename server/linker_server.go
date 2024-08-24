package server

import (
	"context"
	"fmt"
	"time"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/handler"
	"github.com/sicozz/project00/utils"
)

type LinkerService struct {
	proto00.UnimplementedLinkerServer
}

func (s *LinkerService) Info(ctx context.Context, req *proto00.InfoReq) (res *proto00.InfoRes, err error) {
	programInfo := handler.QueryGetServiceInfo()
	return &proto00.InfoRes{
		Version: programInfo.Version,
		Banner:  programInfo.Banner,
	}, nil
}

func (s *LinkerService) Subscribe(req *proto00.SubscribeReq, stream proto00.Linker_SubscribeServer) error {
	term := 1000
	for {
		stream.Send(&proto00.Heartbeat{Term: fmt.Sprintf("%v", term)})
		utils.Info(fmt.Sprintf("Sent heartbeat: %v", term))
		time.Sleep(1 * time.Second)
		term = term + 1
	}
}
