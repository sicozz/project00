package utils

import (
	"math/rand"
	"time"
)

func RandomElectionTimeout() time.Duration {
	randRange := RAFT_MAX_ELECTION_TOUT - RAFT_MIN_ELECTION_TOUT
	randTime := (rand.Int63n(int64(randRange)) + RAFT_MIN_ELECTION_TOUT) * int64(time.Millisecond)
	return time.Duration(randTime)
}
