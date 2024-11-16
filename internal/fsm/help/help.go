package help

import (
	"context"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/fsm"
	"github.com/oke11o/sb-habits-bot/internal/fsm/base"
	"github.com/oke11o/sb-habits-bot/internal/fsm/sender"
)

const InstructionText = `*Инструкция*
Тут можно сохранять заметки
`

const HelpCommand = "/help"

func NewHelp(deps *fsm.Deps) *Help {
	return &Help{
		Base: base.Base{Deps: deps},
	}
}

type Help struct {
	base.Base
}

func (m *Help) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, InstructionText, 0)

	return ctx, smc, state, nil
}
