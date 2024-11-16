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

type HabitSuite struct {
	tests.Suite
	Repo *sqlite.HabitRepo
}

func (s *HabitSuite) SetupSuite() {

}

func (s *HabitSuite) SetupTest() {
	s.InitDb(config.SqliteConfig{
		File:          fmt.Sprintf("../../../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../../../migrations/sqlite",
	}, 111)
	s.Repo = sqlite.NewHabitRepoWithDB(s.DBx)
}

func (s *HabitSuite) TearDownTest() {
	_ = s.DBx.Close()
	//os.Remove(s.DBCfg.File)
}

func (s *HabitSuite) TearDownSuite() {

}

func TestHabitSuite(t *testing.T) {
	suite.Run(t, new(HabitSuite))
}

func (s *HabitSuite) Test_Insert() {
	ctx := context.Background()
	habit := model.Habit{
		UserID:     1,
		Name:       "Утренняя зарядка",
		Type:       "simple",
		Points:     10,
		PointsMode: "fixed",
		CreatedAt:  time.Now(),
	}

	habit, err := s.Repo.CreateHabit(ctx, habit)
	s.Require().NoError(err, "CreateHabit() должен работать без ошибок")
	s.Require().NotZero(habit.ID, "ID привычки должен быть заполнен")

	result, err := s.Repo.GetHabitByID(ctx, habit.ID)
	s.Require().NoError(err, "GetHabitByID() должен работать без ошибок")
	s.Equal(habit.Name, result.Name, "Название привычки должно совпадать")
}

func (s *HabitSuite) Test_Update() {
	ctx := context.Background()
	habit := model.Habit{
		UserID:     1,
		Name:       "Отжимания",
		Type:       "counter",
		Points:     20,
		PointsMode: "proportional",
		CreatedAt:  time.Now(),
	}

	habit, err := s.Repo.CreateHabit(ctx, habit)
	s.Require().NoError(err)

	habit.Points = 30
	habit.Type = "simple"
	err = s.Repo.UpdateHabit(ctx, habit)
	s.Require().NoError(err, "UpdateHabit() должен работать без ошибок")

	updatedHabit, err := s.Repo.GetHabitByID(ctx, habit.ID)
	s.Require().NoError(err)
	s.Equal("simple", updatedHabit.Type, "Тип привычки должен быть обновлён")
	s.Equal(int64(30), updatedHabit.Points, "Баллы должны быть обновлены")
}

func (s *HabitSuite) Test_Upsert() {
	ctx := context.Background()
	habit := model.Habit{
		UserID:     1,
		Name:       "Медитация",
		Type:       "duration",
		Points:     15,
		PointsMode: "fixed",
		CreatedAt:  time.Now(),
	}

	habit, err := s.Repo.UpsertHabit(ctx, habit)
	s.Require().NoError(err)

	habit.Type = "simple"
	habit.Points = 20
	habit, err = s.Repo.UpsertHabit(ctx, habit)
	s.Require().NoError(err)

	upsertedHabit, err := s.Repo.GetHabitByName(ctx, habit.UserID, habit.Name)
	s.Require().NoError(err)
	s.Equal("simple", upsertedHabit.Type)
	s.Equal(int64(20), upsertedHabit.Points)
}

func (s *HabitSuite) Test_DeleteByID() {
	ctx := context.Background()
	habit := model.Habit{
		UserID:    1,
		Name:      "Чтение книги",
		Type:      "simple",
		Points:    5,
		CreatedAt: time.Now(),
	}

	habit, err := s.Repo.CreateHabit(ctx, habit)
	s.Require().NoError(err)

	err = s.Repo.DeleteHabitByID(ctx, habit.ID)
	s.Require().NoError(err)

	_, err = s.Repo.GetHabitByID(ctx, habit.ID)
	s.Require().Error(err, "Привычка должна быть удалена")
}

func (s *HabitSuite) Test_DeleteByName() {
	ctx := context.Background()
	habit := model.Habit{
		UserID:    1,
		Name:      "Физическая активность",
		Type:      "random",
		Points:    10,
		CreatedAt: time.Now(),
	}

	habit, err := s.Repo.CreateHabit(ctx, habit)
	s.Require().NoError(err)

	err = s.Repo.DeleteHabitByName(ctx, habit.UserID, habit.Name)
	s.Require().NoError(err)

	_, err = s.Repo.GetHabitByName(ctx, habit.UserID, habit.Name)
	s.Require().Error(err, "Привычка должна быть удалена")
}

func (s *HabitSuite) Test_GetHabitByID() {
	ctx := context.Background()

	habit := model.Habit{
		UserID:     1,
		Name:       "Зарядка",
		Type:       "simple",
		Points:     10,
		PointsMode: "fixed",
		CreatedAt:  time.Now(),
	}

	habit, err := s.Repo.CreateHabit(ctx, habit)
	s.Require().NoError(err)

	result, err := s.Repo.GetHabitByID(ctx, habit.ID)
	s.Require().NoError(err)
	s.Equal(habit.Name, result.Name, "Название привычки должно совпадать")
	s.Equal(habit.Type, result.Type, "Тип привычки должен совпадать")

	_, err = s.Repo.GetHabitByID(ctx, habit.ID+100)
	s.Require().Error(err, "Ожидается ошибка, если привычка не найдена")
}

func (s *HabitSuite) Test_GetHabitsByUserID() {
	ctx := context.Background()

	habits := []model.Habit{
		{
			UserID:     1,
			Name:       "Зарядка",
			Type:       "simple",
			Points:     10,
			PointsMode: "fixed",
			CreatedAt:  time.Now(),
		},
		{
			UserID:     1,
			Name:       "Отжимания",
			Type:       "counter",
			Points:     20,
			PointsMode: "proportional",
			CreatedAt:  time.Now(),
		},
	}

	for _, habit := range habits {
		_, err := s.Repo.CreateHabit(ctx, habit)
		s.Require().NoError(err)
	}

	result, err := s.Repo.GetHabitsByUserID(ctx, 1)
	s.Require().NoError(err)
	s.Len(result, 2, "Должны быть найдены две привычки")
	s.Equal("Зарядка", result[0].Name)
	s.Equal("Отжимания", result[1].Name)

	emptyResult, err := s.Repo.GetHabitsByUserID(ctx, 2)
	s.Require().NoError(err)
	s.Empty(emptyResult, "Список привычек должен быть пуст, если пользователь не имеет привычек")
}

func (s *HabitSuite) Test_GetHabitByName() {
	ctx := context.Background()

	habit := model.Habit{
		UserID:     1,
		Name:       "Медитация",
		Type:       "duration",
		Points:     15,
		PointsMode: "fixed",
		CreatedAt:  time.Now(),
	}

	habit, err := s.Repo.CreateHabit(ctx, habit)
	s.Require().NoError(err)

	result, err := s.Repo.GetHabitByName(ctx, habit.UserID, habit.Name)
	s.Require().NoError(err)
	s.Equal(habit.Name, result.Name, "Название привычки должно совпадать")
	s.Equal(habit.Type, result.Type, "Тип привычки должен совпадать")

	_, err = s.Repo.GetHabitByName(ctx, habit.UserID, "Неизвестная привычка")
	s.Require().Error(err, "Ожидается ошибка, если привычка не найдена")
}
