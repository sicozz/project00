package statemachine

import (
	"fmt"

	"github.com/sicozz/project00/utils"
)

type StateMachine interface {
	Run() error
}

type RaftSTM struct {
	srvC chan string
}

type state string

const (
	stFollower state = "stFollower"
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

func NewRaftSTM(srvC chan string) *RaftSTM {
	return &RaftSTM{srvC: srvC}
}

func (s *RaftSTM) Run() error {
	go s.handleServerChan()
	go s.handleClientChan()
	return nil
}

func (s *RaftSTM) handleServerChan() error {
	// TODO: Implement events
	for ev := range s.srvC {
		switch ev {
		case string(evHeartbeat):
			utils.Info(fmt.Sprintf("SRV: %v", ev))
		case string(evLeaderTimeout):
			utils.Info(fmt.Sprintf("SRV: %v", ev))
		default:
			utils.Info(fmt.Sprintf("SRV: %v (unknown)", ev))
		}
	}
	return nil
}

func (s *RaftSTM) handleClientChan() error {
	return nil
}
