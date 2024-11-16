package session

import (
	"context"
	"github.com/oke11o/sb-habits-bot/internal/fsm"
	"github.com/oke11o/sb-habits-bot/internal/fsm/sender"
)

func NewSessionMachine(deps *fsm.Deps) *SessionMachine {
	return &SessionMachine{
		deps: deps,
	}
}

type SessionMachine struct {
	deps *fsm.Deps
}

func (s *SessionMachine) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	var scm fsm.Machine
	switch state.Session.Status {
	default:
		scm = sender.NewSenderMachine(s.deps, state.User.ID, "Choose action", 0)
	}

	return ctx, scm, state, nil
}
