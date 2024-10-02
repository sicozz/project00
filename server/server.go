package server

import (
	"context"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/node"
)

type Server interface {
	Shutdown()
}

type Server00 struct {
	stm node.Node
	proto00.UnimplementedLinkerServer
}

func NewServer00(stm node.Node) *Server00 {
	return &Server00{stm: stm}
}

func (s *Server00) Info(
	ctx context.Context,
	req *proto00.InfoReq,
) (*proto00.InfoRes, error) {
	return s.stm.RpcInfo()
}

func (s *Server00) Subscribe(
	req *proto00.SubscribeReq,
	stream proto00.Linker_SubscribeServer,
) error {
	return s.stm.RpcSubscribe(stream)
}

func (s *Server00) RequestVote(
	ctx context.Context,
	req *proto00.RequestVoteReq,
) (*proto00.RequestVoteRes, error) {
	return s.stm.RpcRequestVote()
}

func (s *Server00) Shutdown() {
}
