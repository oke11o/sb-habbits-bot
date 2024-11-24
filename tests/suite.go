package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/oke11o/sb-habits-bot/internal/app"
	"github.com/oke11o/sb-habits-bot/internal/app/migrator"
	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/model/iface"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
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

func (s *Suite) LoadFullConfigFixture(habitRepo iface.HabitRepo, reminderRepo iface.ReminderRepo, userID int64, l *slog.Logger) {
	_, curr, _, ok := runtime.Caller(0)
	s.Require().True(ok)
	cfg, err := s.parseYAML(filepath.Dir(curr) + "/testdata/full_config.yaml")
	s.Require().NoError(err)
	err = app.AddHabitsToDB(context.Background(), habitRepo, reminderRepo, userID, cfg, l)
	s.Require().NoError(err)
}

func (s *Suite) parseYAML(filePath string) (app.Config, error) {
	var config app.Config
	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("os.ReadFile() err: %w", err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("yaml.Unmarshal() err: %w", err)
	}

	return config, nil
}
