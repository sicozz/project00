package statemachine

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/datatype"
	"github.com/sicozz/project00/utils"
	"google.golang.org/grpc"
)

type RaftSTM struct {
	trs           map[event]transition
	st            state
	hosts         map[uuid.UUID]utils.Host
	localhostId   uuid.UUID
	electionTimer *time.Timer
	currentTerm   int
}

type state string

const (
	stCandidate state = "stCandidate"
	stFollower  state = "stFollower"
	stLeader    state = "stLeader"
)

type event string

const (
	evHeartbeat       event = "evHeartbeat"
	evElectionTimeout event = "evElectionTimeout"
)

type transition string

const (
	tsFollowerToLeader transition = "tsFollowerToLeader"
)

func NewRaftSTM(
	hosts map[uuid.UUID]utils.Host,
	selfHostId uuid.UUID,
) *RaftSTM {
	trs := map[event]transition{
		evElectionTimeout: tsFollowerToLeader,
	}

	return &RaftSTM{
		trs:         trs,
		st:          stFollower,
		hosts:       hosts,
		localhostId: selfHostId,
		currentTerm: 1000,
	}
}

func (s *RaftSTM) Run() error {
	go s.handleClient()
	go s.handleElectionTimeout()
	return nil
}

func (s *RaftSTM) handleClient() error {
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

func (s *RaftSTM) handleElectionTimeout() {
	s.resetElectionTimer()
	for {
		<-s.electionTimer.C
		utils.Info("Election timeout, becoming Leader...")
		s.st = stLeader
	}
}

func (s *RaftSTM) resetElectionTimer() {
	if s.electionTimer != nil {
		s.electionTimer.Stop()
	}
	timeout := utils.RandomElectionTimeout()
	s.electionTimer = time.NewTimer(timeout)
}

func (s *RaftSTM) RpcInfo() (*proto00.InfoRes, error) {
	programInfo := datatype.ProgramInfo{
		Version: utils.VERSION,
		Banner:  utils.BANNER,
	}
	return &proto00.InfoRes{
		Version: programInfo.Version,
		Banner:  programInfo.Banner,
		Term:    int32(s.currentTerm),
	}, nil
}

func (s *RaftSTM) RpcSubscribe(
	stream proto00.Linker_SubscribeServer,
) error {
	for {
		select {
		case <-stream.Context().Done():
			err := stream.Context().Err()
			utils.Info(fmt.Sprintf("Follower unsubscribed: %v", err))
			return err
		default:
			switch s.st {
			case stFollower:
				utils.Warn("Follower should redirect")
				return nil
			case stLeader:
				stream.Send(&proto00.Heartbeat{Term: fmt.Sprintf("%v", s.currentTerm)})
				s.currentTerm += 1
				time.Sleep(1 * time.Second)
			}
		}
	}
}
