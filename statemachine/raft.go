package statemachine

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/datatype"
	"github.com/sicozz/project00/requester"
	"github.com/sicozz/project00/utils"
)

type RaftSTM struct {
	trs           map[event]transition
	st            state
	hosts         map[uuid.UUID]utils.Host
	localhostId   uuid.UUID
	electionTimer *time.Timer
	currentTerm   int
	votedFor      uuid.UUID
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
	s.voteFor(s.localhostId)
	voteCount := 0
	s.resetElectionTimer()
	// TODO: use channels and waitgroups to listen for +1 votedCount to avoid
	// shared variable
	for _, h := range s.hosts {
		if h.Id() == s.localhostId {
			continue
		}
		// TODO: implemente RpcRequestVote
		res, err := requester.RpcRequestVote(h, proto00.RequestVoteReq{})
		if err != nil {
			utils.Error(fmt.Sprintf("Failed to request vote: %v", err))
			continue
		}
		if !res.VoteGranted {
			continue
		}
		voteCount += 1
	}
	utils.Debug(fmt.Sprintf("Got %v votes out of %v", voteCount, len(s.hosts)))
	// Send RequestVote RPCs to all other servers
	// If votes received from majority of servers: become leader
	// If AppendEntries RPC received from new leader: convert to follower
	// If <-s.electionTimer.C start new election
}

func (s *RaftSTM) voteFor(candidateId uuid.UUID) {
	s.votedFor = candidateId
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
	switch s.st {
	case stFollower:
		return &proto00.RequestVoteRes{
			Term:        int32(s.currentTerm),
			VoteGranted: true,
		}, nil
	case stCandidate:
		return &proto00.RequestVoteRes{
			Term:        int32(s.currentTerm),
			VoteGranted: false,
		}, nil
	case stLeader:
		return &proto00.RequestVoteRes{
			Term:        int32(s.currentTerm),
			VoteGranted: false,
		}, nil
	default:
		err := fmt.Errorf("Invalid state, declining vote")
		utils.Error(err.Error())
		return nil, err
	}
}
