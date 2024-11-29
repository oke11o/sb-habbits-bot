package fsm

import (
	"context"
	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/service"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/oke11o/sb-habits-bot/internal/model"
	"github.com/oke11o/sb-habits-bot/internal/model/iface"
)

func NewDeps(cfg config.Config, sessionRepo iface.SessionRepo, sender iface.Sender, doneService *service.Done, logger *slog.Logger) *Deps {
	return &Deps{
		Cfg:         cfg,
		SessionRepo: sessionRepo,
		Sender:      sender,
		Logger:      logger,
		DoneService: doneService,
	}

}

type Deps struct {
	SessionRepo iface.SessionRepo
	Sender      iface.Sender
	Logger      *slog.Logger
	Cfg         config.Config
	DoneService *service.Done
}

type State struct {
	User    model.User
	Session model.Session
	Update  tgbotapi.Update
}

type Machine interface {
	Switch(ctx context.Context, state State) (context.Context, Machine, State, error)
}
