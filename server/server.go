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
	node node.Node
	proto00.UnimplementedLinkerServer
}

func NewServer00(node node.Node) *Server00 {
	return &Server00{node: node}
}

func (s *Server00) Info(
	ctx context.Context,
	req *proto00.InfoReq,
) (*proto00.InfoRes, error) {
	return s.node.RpcInfo()
}

func (s *Server00) Subscribe(
	req *proto00.SubscribeReq,
	stream proto00.Linker_SubscribeServer,
) error {
	return s.node.RpcSubscribe(stream)
}

func (s *Server00) RequestVote(
	ctx context.Context,
	req *proto00.RequestVoteReq,
) (*proto00.RequestVoteRes, error) {
	return s.node.RpcRequestVote()
}

func (s *Server00) Shutdown() {
}
