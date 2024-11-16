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

type ProgressSuite struct {
	tests.Suite
	Repo *sqlite.ProgressRepo
}

func (s *ProgressSuite) SetupSuite() {

}

func (s *ProgressSuite) SetupTest() {
	s.InitDb(config.SqliteConfig{
		File:          fmt.Sprintf("../../../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../../../migrations/sqlite",
	}, 111)
	s.Repo = sqlite.NewProgressRepo(s.DBx)
}

func (s *ProgressSuite) TearDownTest() {
	_ = s.DBx.Close()
	//os.Remove(s.DBCfg.File)
}

func (s *ProgressSuite) TearDownSuite() {
}

func (s *ProgressSuite) Test_CreateProgress() {
	ctx := context.Background()
	progress := model.Progress{
		HabitID:          1,
		UserID:           1,
		AccumulatedValue: 5000,
		Target:           10000,
		LastUpdated:      time.Now(),
	}

	createdProgress, err := s.Repo.CreateProgress(ctx, progress)
	s.Require().NoError(err, "CreateProgress() должен работать без ошибок")
	s.Require().NotZero(createdProgress.HabitID, "HabitID должен быть заполнен")
	s.Equal(progress.AccumulatedValue, createdProgress.AccumulatedValue, "Накопленное значение должно совпадать")
}

func (s *ProgressSuite) Test_UpdateProgress() {
	ctx := context.Background()
	progress := model.Progress{
		HabitID:          1,
		UserID:           1,
		AccumulatedValue: 5000,
		Target:           10000,
		LastUpdated:      time.Now(),
	}

	_, err := s.Repo.CreateProgress(ctx, progress)
	s.Require().NoError(err)

	progress.AccumulatedValue = 7500
	err = s.Repo.UpdateProgress(ctx, progress)
	s.Require().NoError(err, "UpdateProgress() должен работать без ошибок")

	updatedProgress, err := s.Repo.GetProgressByHabitID(ctx, progress.HabitID)
	s.Require().NoError(err)
	s.Equal(int64(7500), updatedProgress.AccumulatedValue, "Накопленное значение должно быть обновлено")
}

func (s *ProgressSuite) Test_DeleteProgress() {
	ctx := context.Background()
	progress := model.Progress{
		HabitID:          1,
		UserID:           1,
		AccumulatedValue: 3000,
		Target:           10000,
		LastUpdated:      time.Now(),
	}

	_, err := s.Repo.CreateProgress(ctx, progress)
	s.Require().NoError(err)

	// Удаляем прогресс
	err = s.Repo.DeleteProgress(ctx, progress.HabitID)
	s.Require().NoError(err, "DeleteProgress() должен работать без ошибок")

	// Проверяем, что прогресс удалён
	_, err = s.Repo.GetProgressByHabitID(ctx, progress.HabitID)
	s.Require().Error(err, "Ожидается ошибка, так как прогресс был удалён")
}

// Тест на метод GetProgressByHabitID
func (s *ProgressSuite) Test_GetProgressByHabitID() {
	ctx := context.Background()
	progress := model.Progress{
		HabitID:          1,
		UserID:           1,
		AccumulatedValue: 2000,
		Target:           10000,
		LastUpdated:      time.Now(),
	}

	_, err := s.Repo.CreateProgress(ctx, progress)
	s.Require().NoError(err)

	// Проверяем получение прогресса по HabitID
	result, err := s.Repo.GetProgressByHabitID(ctx, progress.HabitID)
	s.Require().NoError(err, "GetProgressByHabitID() должен работать без ошибок")
	s.Equal(progress.AccumulatedValue, result.AccumulatedValue, "Накопленное значение должно совпадать")
}

// Тест на метод GetProgressByUserID
func (s *ProgressSuite) Test_GetProgressByUserID() {
	ctx := context.Background()
	userID := int64(1)
	progressList := []model.Progress{
		{HabitID: 1, UserID: userID, AccumulatedValue: 1000, Target: 5000, LastUpdated: time.Now()},
		{HabitID: 2, UserID: userID, AccumulatedValue: 3000, Target: 10000, LastUpdated: time.Now()},
	}

	for _, progress := range progressList {
		_, err := s.Repo.CreateProgress(ctx, progress)
		s.Require().NoError(err)
	}

	// Проверяем получение прогресса по UserID
	result, err := s.Repo.GetProgressByUserID(ctx, userID)
	s.Require().NoError(err, "GetProgressByUserID() должен работать без ошибок")
	s.Len(result, 2, "Должны быть найдены два прогресса")
	s.Equal(progressList[0].AccumulatedValue, result[0].AccumulatedValue)
	s.Equal(progressList[1].AccumulatedValue, result[1].AccumulatedValue)
}

func TestProgressSuite(t *testing.T) {
	suite.Run(t, new(ProgressSuite))
}
