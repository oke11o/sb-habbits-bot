package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/model"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
	"github.com/oke11o/sb-habits-bot/pgk/utils/str"
	"github.com/oke11o/sb-habits-bot/tests"
)

type ReminderSuite struct {
	tests.Suite
	Repo *sqlite.ReminderRepo
}

func (s *ReminderSuite) SetupSuite() {

}

func (s *ReminderSuite) SetupTest() {
	s.InitDb(config.SqliteConfig{
		File:          fmt.Sprintf("../../../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../../../migrations/sqlite",
	}, 111)
	s.Repo = sqlite.NewReminderRepo(s.DBx)
}

func (s *ReminderSuite) TearDownTest() {
	_ = s.DBx.Close()
}

func (s *ReminderSuite) TearDownSuite() {
}

func (s *ReminderSuite) Test_CreateReminder() {
	ctx := context.Background()
	reminder := model.Reminder{
		HabitID: 1,
		UserID:  1,
		Time:    "08:00",
		Days:    model.StringSlice{"mon", "tue", "wed"},
	}

	createdReminder, err := s.Repo.CreateReminder(ctx, reminder)
	s.Require().NoError(err, "CreateReminder() должен работать без ошибок")
	s.Require().NotZero(createdReminder.ID, "ID напоминания должен быть заполнен")
	s.Equal(reminder.Time, createdReminder.Time, "Время напоминания должно совпадать")
}

func (s *ReminderSuite) Test_UpdateReminder() {
	ctx := context.Background()
	reminder := model.Reminder{
		HabitID: 1,
		UserID:  1,
		Time:    "08:00",
		Days:    model.StringSlice{"mon", "tue"},
	}

	createdReminder, err := s.Repo.CreateReminder(ctx, reminder)
	s.Require().NoError(err)

	createdReminder.Time = "09:00"
	createdReminder.Days = model.StringSlice{"wed", "thu"}
	err = s.Repo.UpdateReminder(ctx, createdReminder)
	s.Require().NoError(err, "UpdateReminder() должен работать без ошибок")

	updatedReminder, err := s.Repo.GetRemindersByHabitID(ctx, createdReminder.HabitID)
	s.Require().NoError(err)
	s.Require().Len(updatedReminder, 1, "Должно быть одно напоминание")
	s.Equal("09:00", updatedReminder[0].Time, "Время напоминания должно быть обновлено")
	s.Equal(model.StringSlice{"wed", "thu"}, updatedReminder[0].Days, "Дни напоминания должны быть обновлены")
}

func (s *ReminderSuite) Test_DeleteReminder() {
	ctx := context.Background()
	reminder := model.Reminder{
		HabitID: 1,
		UserID:  1,
		Time:    "10:00",
		Days:    model.StringSlice{"fri"},
	}

	createdReminder, err := s.Repo.CreateReminder(ctx, reminder)
	s.Require().NoError(err)

	err = s.Repo.DeleteReminder(ctx, createdReminder.ID)
	s.Require().NoError(err, "DeleteReminder() должен работать без ошибок")

	reminders, err := s.Repo.GetRemindersByHabitID(ctx, reminder.HabitID)
	s.Require().NoError(err)
	s.Empty(reminders, "Список напоминаний должен быть пуст после удаления")
}

func (s *ReminderSuite) Test_GetRemindersByHabitID() {
	ctx := context.Background()
	habitID := int64(1)
	reminders := []model.Reminder{
		{HabitID: habitID, UserID: 1, Time: "08:00", Days: model.StringSlice{"mon", "tue"}},
		{HabitID: habitID, UserID: 1, Time: "09:00", Days: model.StringSlice{"wed", "thu"}},
	}

	for _, reminder := range reminders {
		_, err := s.Repo.CreateReminder(ctx, reminder)
		s.Require().NoError(err)
	}

	result, err := s.Repo.GetRemindersByHabitID(ctx, habitID)
	s.Require().NoError(err, "GetRemindersByHabitID() должен работать без ошибок")
	s.Len(result, 2, "Должны быть найдены два напоминания")
	s.Equal("08:00", result[0].Time)
	s.Equal("09:00", result[1].Time)
}

func (s *ReminderSuite) Test_GetRemindersByUserID() {
	ctx := context.Background()
	userID := int64(1)
	reminders := []model.Reminder{
		{HabitID: 1, UserID: userID, Time: "07:00", Days: model.StringSlice{"mon"}},
		{HabitID: 2, UserID: userID, Time: "08:00", Days: model.StringSlice{"tue"}},
	}

	for _, reminder := range reminders {
		_, err := s.Repo.CreateReminder(ctx, reminder)
		s.Require().NoError(err)
	}

	result, err := s.Repo.GetRemindersByUserID(ctx, userID)
	s.Require().NoError(err, "GetRemindersByUserID() должен работать без ошибок")
	s.Len(result, 2, "Должны быть найдены два напоминания")
	s.Equal("07:00", result[0].Time)
	s.Equal("08:00", result[1].Time)
}

func TestReminderSuite(t *testing.T) {
	suite.Run(t, new(ReminderSuite))
}
