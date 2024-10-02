package requester

import (
	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/network"
)

func RpcRequestVote(h network.Host, reqData proto00.RequestVoteReq) (proto00.RequestVoteRes, error) {
	return proto00.RequestVoteRes{VoteGranted: true}, nil
}
