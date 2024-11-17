package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/oke11o/sb-habits-bot/internal/app"
)

func RunParseYAMLConfigToDB(ctx context.Context, args []string, appname, version string) error {
	ctx, cfg, l, err := prepareCfgAndLog(ctx, appname, version)
	if err != nil {
		return err
	}
	if len(args) < 3 {
		return errors.New("first arg - userID, second arg - yaml path")
	}
	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user id: %s: err: %w", args[0], err)
	}

	err = app.ParseYAMLConfigToDB(ctx, cfg, userID, args[2], l)
	if err != nil {
		l.ErrorContext(ctx, "app.ParseYAMLConfigToDB error", slog.String("error", err.Error()))
		return fmt.Errorf("app.ParseYAMLConfigToDB error: %w", err)
	}
	return nil
}
