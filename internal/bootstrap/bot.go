package bootstrap

import (
	"context"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/app"
	"github.com/oke11o/sb-habits-bot/internal/handler"
	"github.com/oke11o/sb-habits-bot/internal/model/iface"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
	"github.com/oke11o/sb-habits-bot/internal/service"
	"log/slog"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/log"
)

func Run(ctx context.Context, args []string, appname, version string) error {
	ctx, cfg, l, err := prepareCfgAndLog(ctx, appname, version)
	if err != nil {
		return err
	}
	var userRepo iface.UserRepo
	var incomeRepo iface.IncomeRepo
	var sessionRepo iface.SessionRepo
	switch cfg.DBType {
	case sqlite.DBType:
		sqliteDb, err := sqlite.NewDb(cfg.Sqlite)
		if err != nil {
			l.ErrorContext(ctx, "error sqlite.NewUserRepo()", slog.String("error", err.Error()))
			return fmt.Errorf("sqlite.NewUserRepo() err: %w", err)
		}
		userRepo = sqlite.NewUserRepoWithDB(sqliteDb)
		incomeRepo = sqlite.NewIncomeRepoWithDB(sqliteDb)
		sessionRepo = sqlite.NewSessionRepoWithDB(sqliteDb)
	//case mongo.DBType:
	//	userRepo = mongo.NewUserRepo()
	//case pg.DBType:
	//	userRepo = pg.NewUserRepo()
	default:
		l.ErrorContext(ctx, "error sqlite.NewUserRepo()", slog.String("error", err.Error()))
		return fmt.Errorf("unknown db_type %s", cfg.DBType)
	}
	_ = userRepo
	income := service.NewIncomeServce(userRepo, incomeRepo)
	b := app.NewBot(cfg, l, handler.New(cfg, l, income, userRepo, sessionRepo))
	err = b.Run(ctx)
	if err != nil {
		l.ErrorContext(ctx, "error bot.Run()", slog.String("error", err.Error()))
		return fmt.Errorf("bot.Run() err: %w", err)
	}

	return nil
}

func prepareCfgAndLog(ctx context.Context, appname, version string) (context.Context, config.Config, *slog.Logger, error) {
	l := log.New(true, slog.LevelDebug)
	ctx = log.AppendCtx(ctx, slog.String("version", version))
	err := config.InitDotEnv()
	if err != nil {
		l.ErrorContext(ctx, "error config.InitDotEnv()", slog.String("error", err.Error()))
		return ctx, config.Config{}, l, fmt.Errorf("config.InitDotEnv() err: %w", err)
	}

	cfg, err := config.Load(appname)
	if err != nil {
		l.ErrorContext(ctx, "error config.Load()", slog.String("error", err.Error()))
		return ctx, config.Config{}, l, fmt.Errorf("config.Load() err: %w", err)
	}
	return ctx, cfg, l, nil
}
