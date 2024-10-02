package node

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/info"
	"github.com/sicozz/project00/logger"
	"github.com/sicozz/project00/network"
	"github.com/sicozz/project00/requester"
)

type RaftNode struct {
	trs           map[event]transition
	st            state
	hosts         map[uuid.UUID]network.Host
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
	hosts map[uuid.UUID]network.Host,
	selfHostId uuid.UUID,
) *RaftNode {
	trs := map[event]transition{
		evElectionTimeout: tsFollowerToLeader,
	}

	return &RaftNode{
		trs:         trs,
		st:          stFollower,
		hosts:       hosts,
		localhostId: selfHostId,
		currentTerm: 1000,
	}
}

func (s *RaftNode) Run() error {
	go s.handleElectionTimeout()
	return nil
}

func (s *RaftNode) handleElectionTimeout() {
	s.resetElectionTimer()
	for {
		<-s.electionTimer.C
		s.st = stCandidate
		s.startElection()
	}
}

func (s *RaftNode) resetElectionTimer() {
	if s.electionTimer != nil {
		s.electionTimer.Stop()
	}
	timeout := randomElectionTimeout()
	s.electionTimer = time.NewTimer(timeout)
}

func randomElectionTimeout() time.Duration {
	randRange := RAFT_MAX_ELECTION_TOUT - RAFT_MIN_ELECTION_TOUT
	randTime := (rand.Int63n(int64(randRange)) + RAFT_MIN_ELECTION_TOUT) * int64(time.Millisecond)
	return time.Duration(randTime)
}

func (s *RaftNode) startElection() {
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
			logger.Error(fmt.Sprintf("Failed to request vote: %v", err))
			continue
		}
		if !res.VoteGranted {
			continue
		}
		voteCount += 1
	}
	logger.Debug(fmt.Sprintf("Got %v votes out of %v", voteCount, len(s.hosts)))
	// Send RequestVote RPCs to all other servers
	// If votes received from majority of servers: become leader
	// If AppendEntries RPC received from new leader: convert to follower
	// If <-s.electionTimer.C start new election
}

func (s *RaftNode) voteFor(candidateId uuid.UUID) {
	s.votedFor = candidateId
}

func (s *RaftNode) RpcInfo() (*proto00.InfoRes, error) {
	programInfo := info.ProgramInfo{
		Version: info.VERSION,
		Banner:  info.BANNER,
	}
	return &proto00.InfoRes{
		Version: programInfo.Version,
		Banner:  programInfo.Banner,
		Term:    int32(s.currentTerm),
	}, nil
}

func (s *RaftNode) RpcSubscribe(
	stream proto00.Linker_SubscribeServer,
) error {
	for {
		select {
		case <-stream.Context().Done():
			err := stream.Context().Err()
			logger.Info(fmt.Sprintf("Follower unsubscribed: %v", err))
			return err
		default:
			switch s.st {
			case stFollower:
				logger.Warn("Follower should redirect")
				return nil
			case stCandidate:
				stream.Send(&proto00.Heartbeat{Term: fmt.Sprintf("%v", s.currentTerm)})
				s.currentTerm += 1
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (s *RaftNode) RpcRequestVote() (*proto00.RequestVoteRes, error) {
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
		logger.Error(err.Error())
		return nil, err
	}
}
