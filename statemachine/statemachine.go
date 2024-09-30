package statemachine

import proto00 "github.com/sicozz/project00/api/v0.0"

type StateMachine interface {
	Run() error
	RpcInfo() (*proto00.InfoRes, error)
	RpcSubscribe(proto00.Linker_SubscribeServer) error
	RpcRequestVote() (*proto00.RequestVoteRes, error)
}
