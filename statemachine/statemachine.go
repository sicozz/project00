package statemachine

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/utils"
	"google.golang.org/grpc"
)

type StateMachine interface {
	Run() error
}

type RaftSTM struct {
	srvC       chan string
	trs        map[event]transition
	st         state
	hosts      map[uuid.UUID]utils.Host
	selfHostId uuid.UUID
}

type state string

const (
	stFollower state = "stFollower"
	stLeader   state = "stLeader"
)

type event string

const (
	evHeartbeat     event = "evHeartbeat"
	evLeaderTimeout event = "evLeaderTimeout"
)

type transition string

const (
	tsFollowerToLeader transition = "tsFollowerToLeader"
)

func NewRaftSTM(srvC chan string, hosts map[uuid.UUID]utils.Host, selfHostId uuid.UUID) *RaftSTM {
	trs := map[event]transition{
		evLeaderTimeout: tsFollowerToLeader,
	}

	return &RaftSTM{srvC: srvC, trs: trs, st: stFollower, hosts: hosts, selfHostId: selfHostId}
}

func (s *RaftSTM) Run() error {
	go s.handleServerChan()
	go s.handleClientChan()
	return nil
}

func (s *RaftSTM) handleServerChan() error {
	// TODO: Implement events
	for ev := range s.srvC {
		switch s.st {
		case stFollower:
			switch ev {
			case string(evHeartbeat):
				utils.Info(fmt.Sprintf("[Server] Follower: %v", ev))
			case string(evLeaderTimeout):
				utils.Info(fmt.Sprintf("[Server] Follower: %v", ev))
			default:
				utils.Info(fmt.Sprintf("[Server] Follower: %v (unknown)", ev))
			}
		default:
			utils.Error(fmt.Sprintf("Bad stm.st %v", s.st))
		}
	}
	return nil
}

func (s *RaftSTM) handleClientChan() error {
	// create client
	comodin := 0
	targetNode := fmt.Sprintf("node0%v:5005%v", comodin, comodin)
	conn, err := grpc.Dial(targetNode, grpc.WithInsecure())
	for err != nil {
		utils.Error(fmt.Sprintf("[Client] Failed to connect: %v", err))
		comodin = 1 - comodin
		targetNode = fmt.Sprintf("node0%v:5005%v", comodin, comodin)
		conn, err = grpc.Dial(targetNode, grpc.WithInsecure())
		// return err
	}
	defer conn.Close()
	client := proto00.NewLinkerClient(conn)

	switch s.st {
	case stFollower:
		req := &proto00.SubscribeReq{}
		stream, err := client.Subscribe(context.Background(), req)
		if err != nil {
			utils.Error(fmt.Sprintf("[Client] Error while calling Subscribe RPC: %v", err))
			return err
		}
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				utils.Error(fmt.Sprintf("[Client] Error while receiving heartbeat: %v", err))
			}
			utils.Info(fmt.Sprintf("[Client] Heartbeat [%v]", resp))
		}
	default:
		utils.Error(fmt.Sprintf("Bad stm.st %v", s.st))
	}
	return nil
}
