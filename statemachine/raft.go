package statemachine

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/datatype"
	"github.com/sicozz/project00/utils"
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
	go s.handleElectionTimeout()
	return nil
}

func (s *RaftSTM) handleElectionTimeout() {
	s.resetElectionTimer()
	for {
		<-s.electionTimer.C
		utils.Info("Election timeout, becoming Candidate...")
		s.st = stCandidate
		s.startElection()
	}
}

func (s *RaftSTM) resetElectionTimer() {
	if s.electionTimer != nil {
		s.electionTimer.Stop()
	}
	timeout := utils.RandomElectionTimeout()
	s.electionTimer = time.NewTimer(timeout)
}

func (s *RaftSTM) startElection() {
	s.currentTerm += 100
	s.resetElectionTimer()
	// Send RequestVote RPCs to all other servers
	// If votes received from majority of servers: become leader
	// If AppendEntries RPC received from new leader: convert to follower
	// If <-s.electionTimer.C start new election
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
			case stCandidate:
				stream.Send(&proto00.Heartbeat{Term: fmt.Sprintf("%v", s.currentTerm)})
				s.currentTerm += 1
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (s *RaftSTM) RpcRequestVote() (*proto00.RequestVoteRes, error) {
	return &proto00.RequestVoteRes{
		Term:        int32(s.currentTerm),
		VoteGranted: true,
	}, nil
}
