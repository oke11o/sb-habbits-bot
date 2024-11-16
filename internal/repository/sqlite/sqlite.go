package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/oke11o/sb-habits-bot/internal/config"
)

const DBType = "sqlite"

func NewDb(cfg config.SqliteConfig) (*sqlx.DB, error) {
	db, err := sql.Open("sqlite3", cfg.File)
	if err != nil {
		return nil, fmt.Errorf("sql.Open() err: %w", err)
	}
	return sqlx.NewDb(db, "sqlite3"), nil
}
