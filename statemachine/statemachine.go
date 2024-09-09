package statemachine

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/utils"
	"google.golang.org/grpc"
)

type StateMachine interface {
	Run() error
}

type RaftSTM struct {
	srvC        chan string
	trs         map[event]transition
	st          state
	hosts       map[uuid.UUID]utils.Host
	localhostId uuid.UUID
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

func NewRaftSTM(
	srvC chan string,
	hosts map[uuid.UUID]utils.Host,
	selfHostId uuid.UUID,
) *RaftSTM {
	trs := map[event]transition{
		evLeaderTimeout: tsFollowerToLeader,
	}

	return &RaftSTM{
		srvC:        srvC,
		trs:         trs,
		st:          stFollower,
		hosts:       hosts,
		localhostId: selfHostId,
	}
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
		case stLeader:
			switch ev {
			case string(evHeartbeat):
				utils.Info(fmt.Sprintf("[Server] Leader: %v", ev))
			case string(evLeaderTimeout):
				utils.Info(fmt.Sprintf("[Server] Leader: %v", ev))
			default:
				utils.Info(fmt.Sprintf("[Server] Leader: %v (unknown)", ev))
			}
		default:
			utils.Error(fmt.Sprintf("Bad stm.st %v", s.st))
		}
	}
	return nil
}

func (s *RaftSTM) handleClientChan() error {
	switch s.st {
	case stFollower:
		targetNode := "192.168.100.10:50050"
		conn, err := grpc.Dial(targetNode, grpc.WithInsecure())
		if err != nil {
			utils.Error(fmt.Sprintf("[Client] Failed to connect: %v", err))
			conn, err = grpc.Dial(targetNode, grpc.WithInsecure())
			// return err
		}
		defer conn.Close()
		client := proto00.NewLinkerClient(conn)
		req := &proto00.SubscribeReq{}
		stream, err := client.Subscribe(context.Background(), req)
		if err != nil {
			utils.Error(
				fmt.Sprintf(
					"[Client] Error while calling Subscribe RPC: %v",
					err,
				),
			)
			return err
		}
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				utils.Error(
					fmt.Sprintf(
						"[Client] Error while receiving heartbeat: %v",
						err,
					),
				)
			}
			utils.Info(fmt.Sprintf("[Client] Heartbeat [%v]", resp))
		}
	case stLeader:
		time.Sleep(3 * time.Second)
		targetNode := "192.168.100.11:50051"
		conn, err := grpc.Dial(targetNode, grpc.WithInsecure())
		if err != nil {
			utils.Error(fmt.Sprintf("[Client] Failed to connect: %v", err))
			return err
		}
		defer conn.Close()
		client := proto00.NewLinkerClient(conn)
		res, err := client.Info(context.Background(), &proto00.InfoReq{})
		if err != nil {
			utils.Error(
				fmt.Sprintf(
					"[Client] Error while calling Info RPC: %v",
					err,
				),
			)
			return err
		}
		utils.Debug(fmt.Sprintf("INFO RES:\t%v", res))
	default:
		utils.Error(fmt.Sprintf("Bad stm.st %v", s.st))
	}
	return nil
}

func (s *RaftSTM) TmpSetLeaderState() {
	s.st = stLeader
}
