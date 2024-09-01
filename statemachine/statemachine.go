package statemachine

import (
	"fmt"

	"github.com/sicozz/project00/utils"
)

type StateMachine interface {
	Run() error
}

type StateMachine00 struct {
	srvC chan string
}

func NewStateMachine00(srvC chan string) *StateMachine00 {
	return &StateMachine00{srvC: srvC}
}

func (s *StateMachine00) Run() error {
	for n := range s.srvC {
		utils.Info(fmt.Sprintf("STM (from: SRV): Unsubscribed at %v", n))
	}
	return nil
}
