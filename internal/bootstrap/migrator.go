package bootstrap

import (
	"context"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/app/migrator"
	"log/slog"
)

func RunMigrator(ctx context.Context, _ []string, appname, version string) error {
	ctx, cfg, l, err := prepareCfgAndLog(ctx, appname, version)
	if err != nil {
		return err
	}

	err = migrator.RunMigrator(ctx, cfg.Sqlite)
	if err != nil {
		l.ErrorContext(ctx, "app.RunMigrator error", slog.String("error", err.Error()))
		return fmt.Errorf("app.RunMigrator error: %w", err)
	}
	return nil
}
