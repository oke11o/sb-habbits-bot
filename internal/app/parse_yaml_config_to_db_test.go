package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
	"github.com/oke11o/sb-habits-bot/tests"
)

type AddHabitsSuite struct {
	tests.Suite
	HabitRepo    *sqlite.HabitRepo
	ReminderRepo *sqlite.ReminderRepo
	Logger       *slog.Logger
}

func (s *AddHabitsSuite) SetupTest() {
	s.InitDb(config.SqliteConfig{
		File:          fmt.Sprintf("../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../migrations/sqlite",
	}, 111)
	s.HabitRepo = sqlite.NewHabitRepoWithDB(s.DBx)
	s.ReminderRepo = sqlite.NewReminderRepo(s.DBx)
	s.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func (s *AddHabitsSuite) TearDownTest() {
	_ = s.DBx.Close()
}

func (s *AddHabitsSuite) Test_AddHabitsToDB() {
	ctx := context.Background()
	userID := int64(1)

	// Определяем конфигурацию для теста
	config := Config{
		Habits: []HabitConfig{
			{
				Name:   "Утренняя зарядка",
				Type:   "simple",
				Points: 10,
				Reminder: struct {
					Time string   `yaml:"time"`
					Days []string `yaml:"days"`
				}{
					Time: "08:00",
					Days: []string{"mon", "tue"},
				},
			},
			{
				Name:       "Отжимания",
				Type:       "counter",
				Target:     30,
				Points:     20,
				PointsMode: "proportional",
			},
		},
	}

	// Запускаем тестируемую функцию
	err := addHabitsToDB(ctx, s.HabitRepo, s.ReminderRepo, userID, config, s.Logger)
	s.Require().NoError(err, "addHabitsToDB() должна завершаться без ошибок")

	// Проверяем добавление первой привычки
	habit, err := s.HabitRepo.GetHabitByName(ctx, userID, "Утренняя зарядка")
	s.Require().NoError(err, "GetHabitByName() должна завершаться без ошибок")
	s.Equal("simple", habit.Type, "Тип привычки должен совпадать")
	s.Equal(int64(10), habit.Points, "Баллы за привычку должны совпадать")

	// Проверяем напоминание для первой привычки
	reminders, err := s.ReminderRepo.GetRemindersByHabitID(ctx, habit.ID)
	s.Require().NoError(err, "GetRemindersByHabitID() должна завершаться без ошибок")
	s.Require().Len(reminders, 1, "Должно быть одно напоминание")
	s.Equal("08:00", reminders[0].Time, "Время напоминания должно совпадать")
	s.Equal("[mon tue]", reminders[0].Days, "Дни напоминания должны совпадать") //TODO: need fix [mon tue]

	// Проверяем добавление второй привычки
	habit, err = s.HabitRepo.GetHabitByName(ctx, userID, "Отжимания")
	s.Require().NoError(err)
	s.Equal("counter", habit.Type, "Тип привычки должен совпадать")
	s.Equal(int64(20), habit.Points, "Баллы за привычку должны совпадать")
	s.Equal("proportional", habit.PointsMode, "Режим начисления баллов должен совпадать")
}

func TestAddHabitsSuite(t *testing.T) {
	suite.Run(t, new(AddHabitsSuite))
}
