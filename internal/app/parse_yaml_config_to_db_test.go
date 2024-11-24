package app

import (
	"context"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/model"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"testing"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
	"github.com/oke11o/sb-habits-bot/tests"
	"github.com/stretchr/testify/suite"
)

func TestAddHabitsSuite(t *testing.T) {
	suite.Run(t, new(AddHabitsSuite))
}

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
	file, err := os.ReadFile("../../tests/testdata/full_config.yaml")
	s.Require().NoError(err, "Ошибка чтения YAML файла")

	cfg := Config{}
	err = yaml.Unmarshal(file, &cfg)
	s.Require().NoError(err, "Ошибка парсинга YAML файла")
	ctx := context.Background()
	userID := int64(1)

	err = AddHabitsToDB(ctx, s.HabitRepo, s.ReminderRepo, userID, cfg, s.Logger)
	s.Require().NoError(err, "AddHabitsToDB() должна завершаться без ошибок")

	for _, habitConfig := range cfg.Habits {
		habit, err := s.HabitRepo.GetHabitByName(ctx, userID, habitConfig.Name)
		s.Require().NoError(err, "GetHabitByName() должна завершаться без ошибок")
		s.Equal(habitConfig.Type, habit.Type, "Тип привычки должен совпадать")
		s.Equal(habitConfig.Points, habit.Points, "Баллы за привычку должны совпадать")

		if habitConfig.Reminder.Time != "" {
			reminders, err := s.ReminderRepo.GetRemindersByHabitID(ctx, habit.ID)
			s.Require().NoError(err, "GetRemindersByHabitID() должна завершаться без ошибок")
			s.Require().Len(reminders, 1, "Должно быть одно напоминание")
			s.Equal(habitConfig.Reminder.Time, reminders[0].Time, "Время напоминания должно совпадать")
			s.Equal(habitConfig.Reminder.Days, []string(reminders[0].Days), "Дни напоминания должны совпадать")
		}
	}

	// Проверка специфичных полей для некоторых типов привычек
	timeHabit, err := s.HabitRepo.GetHabitByName(ctx, userID, "Проснуться в 6 утра")
	s.Require().NoError(err)
	s.Equal("06:00", timeHabit.TargetTime)
	s.Equal("08:00", timeHabit.MaxTime)

	durationHabit, err := s.HabitRepo.GetHabitByName(ctx, userID, "Медитация")
	s.Require().NoError(err)
	s.Equal("15m", durationHabit.TargetDuration)

	cumulativeHabit, err := s.HabitRepo.GetHabitByName(ctx, userID, "Шаги в неделю")
	s.Require().NoError(err)
	s.Equal("steps", cumulativeHabit.Unit)
	s.Equal(int64(70000), cumulativeHabit.Target)

	periodicHabit, err := s.HabitRepo.GetHabitByName(ctx, userID, "Бег по утрам")
	s.Require().NoError(err)
	s.Equal(int64(2), periodicHabit.IntervalDays)

	checklistHabit, err := s.HabitRepo.GetHabitByName(ctx, userID, "Утренний ритуал")
	s.Require().NoError(err)
	s.Equal("checklist", checklistHabit.Type)
	s.Equal(model.StringSlice{"Зарядка", "Душ", "Завтрак"}, checklistHabit.Tasks)

	randomHabit, err := s.HabitRepo.GetHabitByName(ctx, userID, "Физическая активность")
	s.Require().NoError(err)
	s.Equal(model.StringSlice{"Бег", "Отжимания", "Йога"}, randomHabit.Options)

	counterHabitExtended, err := s.HabitRepo.GetHabitByName(ctx, userID, "Прогулка на свежем воздухе")
	s.Require().NoError(err)
	s.Equal("steps", counterHabitExtended.Unit)
	s.Equal(int64(5000), counterHabitExtended.Target)
}
