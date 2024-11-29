package done

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

const DoneCommand = "/done"

func NewDone(deps *fsm.Deps) *Done {
	return &Done{
		Base: base.Base{Deps: deps},
	}
}

type Done struct {
	base.Base
}

func (m *Done) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}

	err := m.Deps.DoneService.Done(ctx, state.User.ID, state.Update.Message.Text)
	if err != nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "Успешно добавил запис о вашей привычке", 0)

	return ctx, smc, state, nil
}
