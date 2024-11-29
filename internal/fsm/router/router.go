package router

import (
	"context"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/fsm/done"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/oke11o/sb-habits-bot/internal/fsm"
	"github.com/oke11o/sb-habits-bot/internal/fsm/help"
	"github.com/oke11o/sb-habits-bot/internal/fsm/session"
	"github.com/oke11o/sb-habits-bot/internal/model"
)

const StartCommand = "/start"

type Router struct {
	deps *fsm.Deps
}

func NewRouter(deps *fsm.Deps) (*Router, error) {
	return &Router{
		deps: deps,
	}, nil
}

func (r *Router) GetMachine(ctx context.Context, user model.User, update tgbotapi.Update) (fsm.Machine, fsm.State, error) {
	state := fsm.State{
		User:   user,
		Update: update,
	}
	ses, err := r.deps.SessionRepo.GetOpenedSession(ctx, user.ID)
	if err != nil {
		return nil, state, fmt.Errorf("repo.GetOpenedSession(%d) err: %w", user.ID, err)
	}
	state.Session = ses

	// Получить текущую сессию из базы, если есть
	if update.Message != nil {
		cmdMachine := r.resolveCommandMachine(update)
		if cmdMachine != nil {
			return cmdMachine, state, nil
		}

		return session.NewSessionMachine(r.deps), state, nil
	}

	return nil, state, fmt.Errorf("unknown state machine")
}

func (r *Router) resolveCommandMachine(update tgbotapi.Update) fsm.Machine {

	if strings.HasPrefix(update.Message.Text, help.HelpCommand) || strings.HasPrefix(update.Message.Text, StartCommand) {
		return help.NewHelp(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, done.DoneCommand) {
		return done.NewDone(r.deps)
	}
	return nil
}
