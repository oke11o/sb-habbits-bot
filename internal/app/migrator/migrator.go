package migrator

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/oke11o/sb-habits-bot/internal/config"
)

func RunMigrator(_ context.Context, cfg config.SqliteConfig) error {
	m, err := migrate.New("file://"+cfg.MigrationPath, "sqlite3://"+cfg.File)
	if err != nil {
		return fmt.Errorf("cant migrate.New, err: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate, err: %w", err)
	}
	return nil
}
