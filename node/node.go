package node

import (
	proto00 "github.com/sicozz/project00/api/v0.0"
)

const (
	RAFT_MIN_ELECTION_TOUT = 150
	RAFT_MAX_ELECTION_TOUT = 300
)

type Node interface {
	Run() error
	RpcInfo() (*proto00.InfoRes, error)
	RpcSubscribe(proto00.Linker_SubscribeServer) error
	RpcRequestVote() (*proto00.RequestVoteRes, error)
}
