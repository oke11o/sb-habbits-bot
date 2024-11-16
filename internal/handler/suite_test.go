package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/app/migrator"
	"log/slog"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
	"github.com/oke11o/sb-habits-bot/internal/service"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
)

type Suite struct {
	suite.Suite
	dbx         *sqlx.DB
	userRepo    *sqlite.UserRepo
	sessionRepo *sqlite.SessionRepo
	incomeRepo  *sqlite.IncomeRepo
	dbCfg       config.SqliteConfig
	cfg         config.Config
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {}

func (s *Suite) SetupTest() {
	s.cfg.MaintainerChatID = 111
	s.dbCfg = config.SqliteConfig{
		File:          fmt.Sprintf("../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../migrations/sqlite",
	}
	err := migrator.RunMigrator(context.Background(), s.dbCfg)
	s.Require().NoError(err)

	dbx, err := s.createDB(s.dbCfg)
	s.Require().NoError(err)
	s.dbx = dbx
	s.userRepo = sqlite.NewUserRepoWithDB(dbx)
	s.incomeRepo = sqlite.NewIncomeRepoWithDB(dbx)
	s.sessionRepo = sqlite.NewSessionRepoWithDB(dbx)
	s.T().Logf("Start testing with db in file `%s`", s.dbCfg.File)
}

func (s *Suite) TearDownTest() {
	s.dbx.Close()
	//os.Remove(s.dbCfg.File)
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

func (s *Suite) createHandler() *Handler {
	income := service.NewIncomeServce(s.userRepo, s.incomeRepo)
	h := &Handler{
		logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
		income:      income,
		userRepo:    s.userRepo,
		incomeRepo:  s.incomeRepo,
		sessionRepo: s.sessionRepo,
		cfg:         s.cfg,
	}
	return h
}

type testSender struct {
	assert    func(c tgbotapi.Chattable)
	returnMsg tgbotapi.Message
	returnErr error
}

func (t testSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if t.assert != nil {
		t.assert(c)
	}
	return t.returnMsg, t.returnErr
}
