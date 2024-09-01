package statemachine

import (
	"fmt"

	"github.com/sicozz/project00/utils"
)

type StateMachine00 struct {
	SRVCh chan string
}

func NewStateMachine00(srvCh chan string) *StateMachine00 {
	return &StateMachine00{SRVCh: srvCh}
}

func (s *StateMachine00) Run() error {
	n := <-s.SRVCh
	utils.Info(fmt.Sprintf("SRVCH: ON %v", n))
	return nil
}
