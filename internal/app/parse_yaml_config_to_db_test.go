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

	// Определяем расширенную конфигурацию для теста
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
			{
				Name:       "Проснуться в 6 утра",
				Type:       "time",
				TargetTime: "06:00",
				MaxTime:    "08:00",
				Points:     100,
				PointsMode: "time_based",
			},
			{
				Name:           "Медитация",
				Type:           "duration",
				TargetDuration: "15m",
				Points:         15,
				PointsMode:     "proportional",
			},
			{
				Name:       "Шаги в неделю",
				Type:       "cumulative",
				Target:     70000,
				Unit:       "steps",
				Points:     50,
				PointsMode: "proportional",
			},
			{
				Name:         "Бег по утрам",
				Type:         "periodic",
				IntervalDays: 2,
				Points:       15,
			},
			{
				Name:   "Утренний ритуал",
				Type:   "checklist",
				Tasks:  []string{"Зарядка", "Душ", "Завтрак"},
				Points: 20,
			},
			{
				Name:    "Физическая активность",
				Type:    "random",
				Options: []string{"Бег", "Отжимания", "Йога"},
				Points:  10,
			},
		},
	}

	// Запускаем тестируемую функцию
	err := addHabitsToDB(ctx, s.HabitRepo, s.ReminderRepo, userID, config, s.Logger)
	s.Require().NoError(err, "addHabitsToDB() должна завершаться без ошибок")

	// Проверяем все добавленные привычки
	for _, habitConfig := range config.Habits {
		habit, err := s.HabitRepo.GetHabitByName(ctx, userID, habitConfig.Name)
		s.Require().NoError(err, "GetHabitByName() должна завершаться без ошибок")
		s.Equal(habitConfig.Type, habit.Type, "Тип привычки должен совпадать")
		s.Equal(habitConfig.Points, habit.Points, "Баллы за привычку должны совпадать")

		// Проверка напоминания, если оно задано
		if habitConfig.Reminder.Time != "" {
			reminders, err := s.ReminderRepo.GetRemindersByHabitID(ctx, habit.ID)
			s.Require().NoError(err, "GetRemindersByHabitID() должна завершаться без ошибок")
			s.Require().Len(reminders, 1, "Должно быть одно напоминание")
			s.Equal(habitConfig.Reminder.Time, reminders[0].Time, "Время напоминания должно совпадать")
			s.Equal(fmt.Sprintf("%v", habitConfig.Reminder.Days), reminders[0].Days, "Дни напоминания должны совпадать")
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
	s.Equal([]string{"Зарядка", "Душ", "Завтрак"}, checklistHabit.Tasks)

	randomHabit, err := s.HabitRepo.GetHabitByName(ctx, userID, "Физическая активность")
	s.Require().NoError(err)
	s.Equal([]string{"Бег", "Отжимания", "Йога"}, randomHabit.Options)
}

func TestAddHabitsSuite(t *testing.T) {
	suite.Run(t, new(AddHabitsSuite))
}
