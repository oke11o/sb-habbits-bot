package handler

import (
	"context"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
	"github.com/oke11o/sb-habits-bot/internal/service"
	"log/slog"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/fsm"
	"github.com/oke11o/sb-habits-bot/internal/fsm/router"
	"github.com/oke11o/sb-habits-bot/internal/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/oke11o/sb-habits-bot/internal/model"
	"github.com/oke11o/sb-habits-bot/internal/model/iface"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
)

type incomer interface {
	Income(ctx context.Context, requestID string, update tgbotapi.Update) (model.User, error)
}

func New(cfg config.Config, l *slog.Logger, income incomer, userRepo iface.UserRepo, sessionRepo iface.SessionRepo, doneService *service.Done) *Handler {
	return &Handler{
		cfg:         cfg,
		logger:      l,
		income:      income,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		doneService: doneService,
	}
}

type Handler struct {
	cfg         config.Config
	logger      *slog.Logger
	sender      iface.Sender
	income      incomer
	userRepo    iface.UserRepo
	sessionRepo iface.SessionRepo
	incomeRepo  *sqlite.IncomeRepo
	doneService *service.Done
}

func (h *Handler) SetSender(sender iface.Sender) {
	h.sender = sender
}

func (h *Handler) HandleUpdate(ctx context.Context, update tgbotapi.Update) error {
	requestID := fmt.Sprintf("%s-%d", str.RandStringRunes(32, ""), update.UpdateID)
	ctx = log.AppendCtx(ctx, slog.String("request_id", requestID))

	user, err := h.income.Income(ctx, requestID, update)
	if err != nil {
		h.logger.ErrorContext(ctx, "income.Income Error", "error", err.Error())
		return fmt.Errorf("income.Income() err: %w", err)
	}
	ctx = log.AppendCtx(ctx, slog.Int64("user_id", user.ID))

	deps := fsm.NewDeps(h.cfg, h.sessionRepo, h.sender, h.doneService, h.logger)
	routr, err := router.NewRouter(deps)
	if err != nil {
		h.logger.ErrorContext(ctx, "fsm.NewRouter Error", "error", err.Error())
		return fmt.Errorf("fsm.NewRouter() err: %w", err)
	}
	machine, state, err := routr.GetMachine(ctx, user, update)
	if err != nil {
		h.logger.ErrorContext(ctx, "router.GetMachine Error", "error", err.Error())
		return fmt.Errorf("router.GetMachine() err: %w", err)
	}

	for machine != nil {
		ctx, machine, state, err = machine.Switch(ctx, state)
		if err != nil {
			h.logger.ErrorContext(ctx, "machine.Switch Error", "error", err.Error())
			return fmt.Errorf("machine.Switch() err: %w", err)
		}
	}

	return nil
}
