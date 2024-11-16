package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/oke11o/sb-habits-bot/internal/app/migrator"
	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DBx   *sqlx.DB
	DBCfg config.SqliteConfig
	Cfg   config.Config
}

func (s *Suite) SetupSuite() {}

func (s *Suite) SetupTest() {
	s.InitDb(config.SqliteConfig{
		File:          fmt.Sprintf("../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../migrations/sqlite",
	}, 111)
}

func (s *Suite) TearDownTest() {
	s.DBx.Close()
	//os.Remove(s.DBCfg.File)
}

func (s *Suite) TearDownSuite() {}

func (s *Suite) createDB(cfg config.SqliteConfig) (*sqlx.DB, error) {
	db, err := sql.Open("sqlite3", cfg.File)
	if err != nil {
		return nil, fmt.Errorf("sql.Open() err: %w", err)
	}
	dbx := sqlx.NewDb(db, "sqlite3")
	return dbx, nil
}

func (s *Suite) InitDb(cfg config.SqliteConfig, maintainerChatID int64) {
	s.Cfg.MaintainerChatID = maintainerChatID
	s.DBCfg = cfg
	err := migrator.RunMigrator(context.Background(), s.DBCfg)
	s.Require().NoError(err)

	dbx, err := s.createDB(s.DBCfg)
	s.Require().NoError(err)
	s.DBx = dbx
	s.T().Logf("Start testing with db in file `%s`", s.DBCfg.File)
}
