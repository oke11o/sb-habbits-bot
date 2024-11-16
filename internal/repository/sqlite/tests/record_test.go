package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/model"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
	"github.com/oke11o/sb-habits-bot/tests"
)

type RecordSuite struct {
	tests.Suite
	Repo *sqlite.RecordRepo
}

func (s *RecordSuite) SetupSuite() {

}

func (s *RecordSuite) SetupTest() {
	s.InitDb(config.SqliteConfig{
		File:          fmt.Sprintf("../../../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../../../migrations/sqlite",
	}, 111)
	s.Repo = sqlite.NewRecordRepo(s.DBx)
}

func (s *RecordSuite) TearDownTest() {
	s.DBx.Close()
	//os.Remove(s.DBCfg.File)
}

func (s *RecordSuite) TearDownSuite() {
}

func (s *RecordSuite) Test_CreateRecord() {
	ctx := context.Background()
	record := model.Record{
		HabitID:   1,
		UserID:    1,
		Value:     10,
		Timestamp: time.Now(),
		Points:    5,
	}

	createdRecord, err := s.Repo.CreateRecord(ctx, record)
	s.Require().NoError(err, "CreateRecord() должен работать без ошибок")
	s.Require().NotZero(createdRecord.ID, "ID записи должен быть заполнен")
	s.Equal(record.Value, createdRecord.Value, "Значение записи должно совпадать")
}

func (s *RecordSuite) Test_DeleteRecord() {
	ctx := context.Background()
	record := model.Record{
		HabitID:   1,
		UserID:    1,
		Value:     15,
		Timestamp: time.Now(),
		Points:    10,
	}

	record, err := s.Repo.CreateRecord(ctx, record)
	s.Require().NoError(err)

	err = s.Repo.DeleteRecord(ctx, record.ID)
	s.Require().NoError(err, "DeleteRecord() должен работать без ошибок")

	_, err = s.Repo.GetLatestRecordByHabitID(ctx, record.HabitID)
	s.Require().Error(err, "Ожидается ошибка, так как запись была удалена")
}

func (s *RecordSuite) Test_GetRecordsByHabitID() {
	ctx := context.Background()
	habitID := int64(1)
	records := []model.Record{
		{HabitID: habitID, UserID: 1, Value: 10, Timestamp: time.Now(), Points: 5},
		{HabitID: habitID, UserID: 1, Value: 20, Timestamp: time.Now().Add(time.Second), Points: 10},
	}

	for _, record := range records {
		_, err := s.Repo.CreateRecord(ctx, record)
		s.Require().NoError(err)
	}

	result, err := s.Repo.GetRecordsByHabitID(ctx, habitID, 10)
	s.Require().NoError(err, "GetRecordsByHabitID() должен работать без ошибок")
	s.Len(result, 2, "Должно быть возвращено две записи")
	s.Equal(records[1].Value, result[0].Value, "Последняя запись должна быть первой в списке")
	s.Equal(records[0].Value, result[1].Value, "Первая запись должна быть последней в списке")
}

func (s *RecordSuite) Test_GetRecordsByUserID() {
	ctx := context.Background()
	userID := int64(1)
	records := []model.Record{
		{HabitID: 1, UserID: userID, Value: 10, Timestamp: time.Now(), Points: 5},
		{HabitID: 2, UserID: userID, Value: 30, Timestamp: time.Now().Add(time.Second), Points: 15},
	}

	for _, record := range records {
		_, err := s.Repo.CreateRecord(ctx, record)
		s.Require().NoError(err)
	}

	result, err := s.Repo.GetRecordsByUserID(ctx, userID, 10)
	s.Require().NoError(err, "GetRecordsByUserID() должен работать без ошибок")
	s.Len(result, 2, "Должно быть возвращено две записи")
	s.Equal(records[1].Points, result[0].Points, "Последняя запись должна быть первой в списке")
	s.Equal(records[0].Points, result[1].Points, "Первая запись должна быть последней в списке")
}

func (s *RecordSuite) Test_GetLatestRecordByHabitID() {
	ctx := context.Background()
	habitID := int64(1)
	records := []model.Record{
		{HabitID: habitID, UserID: 1, Value: 10, Timestamp: time.Now().Add(-time.Hour), Points: 5},
		{HabitID: habitID, UserID: 1, Value: 20, Timestamp: time.Now(), Points: 10},
	}

	for _, record := range records {
		_, err := s.Repo.CreateRecord(ctx, record)
		s.Require().NoError(err)
	}

	latestRecord, err := s.Repo.GetLatestRecordByHabitID(ctx, habitID)
	s.Require().NoError(err, "GetLatestRecordByHabitID() должен работать без ошибок")
	s.Equal(records[1].Value, latestRecord.Value, "Последняя запись должна быть возвращена")
}

func TestRecordSuite(t *testing.T) {
	suite.Run(t, new(RecordSuite))
}
