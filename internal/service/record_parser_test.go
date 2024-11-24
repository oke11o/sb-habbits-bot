package service

import (
	"context"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/model"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
	"github.com/oke11o/sb-habits-bot/tests"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestAddHabitsSuite(t *testing.T) {
	suite.Run(t, new(RecordParserSuite))
}

type RecordParserSuite struct {
	tests.Suite
	HabitRepo    *sqlite.HabitRepo
	ReminderRepo *sqlite.ReminderRepo
	Logger       *slog.Logger
	userID       int64
}

func (s *RecordParserSuite) SetupTest() {
	s.InitDb(config.SqliteConfig{
		File:          fmt.Sprintf("../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../migrations/sqlite",
	}, 111)
	s.HabitRepo = sqlite.NewHabitRepoWithDB(s.DBx)
	s.ReminderRepo = sqlite.NewReminderRepo(s.DBx)
	s.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	s.userID = 1

	s.LoadFullConfigFixture(s.HabitRepo, s.ReminderRepo, s.userID, s.Logger)
}

func (s *RecordParserSuite) TearDownTest() {
	_ = s.DBx.Close()
}

func (s *RecordParserSuite) TestParseCommand() {
	nowTime, err := time.Parse(time.RFC3339, "2024-01-02T15:04:05+00:00")
	s.Require().NoError(err)

	now = func() time.Time {
		return nowTime
	}
	expect := map[string]model.Record{
		"/done Утренняя зарядка": {
			Timestamp: nowTime,
			HabitID:   1,
			UserID:    s.userID,
			Value:     1,
			Points:    10,
		},
		"/done Отжимания 99": {
			Timestamp: nowTime,
			HabitID:   2,
			UserID:    s.userID,
			Value:     99,
			Points:    20,
		},
		"/done Отжимания 15": {
			Timestamp: nowTime,
			HabitID:   2,
			UserID:    s.userID,
			Value:     15,
			Points:    10,
		},
		"/done Прогулка на свежем воздухе 4900": {
			Timestamp: nowTime,
			HabitID:   3,
			UserID:    s.userID,
			Value:     4900,
			Points:    39,
		},
		"/done Проснуться в 6 утра 07:05": {
			Timestamp: nowTime,
			HabitID:   4,
			UserID:    s.userID,
			Value:     1,
			Points:    45,
		},
		"/done Медитация 15m": {
			Timestamp: nowTime,
			HabitID:   5,
			UserID:    s.userID,
			Value:     0,
			Points:    0,
		},
		"/done Шаги в неделю 1000": {
			Timestamp: nowTime,
			HabitID:   6,
			UserID:    s.userID,
			Value:     0,
			Points:    0,
		},
		"/done Шаги в неделю 2000": {
			Timestamp: nowTime,
			HabitID:   6,
			UserID:    s.userID,
			Value:     0,
			Points:    0,
		},
		"/done Бег по утрам": {
			Timestamp: nowTime,
			HabitID:   7,
			UserID:    s.userID,
			Value:     0,
			Points:    0,
		},
		"/done Утренний ритуал": {
			Timestamp: nowTime,
			HabitID:   8,
			UserID:    s.userID,
			Value:     0,
			Points:    0,
		},
		"/done Физическая активность": {
			Timestamp: nowTime,
			HabitID:   9,
			UserID:    s.userID,
			Value:     0,
			Points:    0,
		},
	}
	msgs := []string{
		"/done Утренняя зарядка",
		"/done Отжимания 99",
		"/done Отжимания 15",
		"/done Прогулка на свежем воздухе 4900",
		"/done Проснуться в 6 утра 07:05",
		//"/done Медитация 15m",
		//"/done Шаги в неделю 1000",
		//"/done Шаги в неделю 2000",
		//"/done Бег по утрам",
		//"/done Утренний ритуал",
		//"/done Физическая активность",
	}

	parser := NewRecordParser(s.HabitRepo)
	for _, msg := range msgs {
		ctx := context.Background()
		record, err := parser.ParseCommand(ctx, 1, msg)
		s.Require().NoError(err)
		s.Require().Equal(expect[msg], record)

	}
}
