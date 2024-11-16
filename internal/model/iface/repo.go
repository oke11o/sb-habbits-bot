package iface

import (
	"context"

	"github.com/oke11o/sb-habits-bot/internal/model"
)

type IncomeRepo interface {
	SaveIncome(ctx context.Context, income model.IncomeRequest) (model.IncomeRequest, error)
}
type SessionRepo interface {
	SaveSession(ctx context.Context, session model.Session) (model.Session, error)
	GetOpenedSession(ctx context.Context, userID int64) (model.Session, error)
	CloseSession(ctx context.Context, session model.Session) error
}

type UserRepo interface {
	SaveUser(ctx context.Context, user model.User) (model.User, error)
	UserCount(ctx context.Context) (int, error)
}

type HabitRepo interface {
	CreateHabit(ctx context.Context, habit model.Habit) (model.Habit, error)
	UpdateHabit(ctx context.Context, habit model.Habit) error
	UpsertHabit(ctx context.Context, habit model.Habit) (model.Habit, error)

	GetHabitByID(ctx context.Context, habitID int64) (model.Habit, error)
	GetHabitsByUserID(ctx context.Context, userID int64) ([]model.Habit, error)
	GetHabitByName(ctx context.Context, userID int64, name string) (model.Habit, error)

	DeleteHabitByID(ctx context.Context, habitID int64) error
	DeleteHabitByName(ctx context.Context, userID int64, name string) error
}

type RecordRepo interface {
	CreateRecord(ctx context.Context, record model.Record) (model.Record, error)
	DeleteRecord(ctx context.Context, recordID int64) error

	GetRecordsByHabitID(ctx context.Context, habitID int64, limit int) ([]model.Record, error)
	GetRecordsByUserID(ctx context.Context, userID int64, limit int) ([]model.Record, error)
	GetLatestRecordByHabitID(ctx context.Context, habitID int64) (model.Record, error)
}

type ReminderRepo interface {
	CreateReminder(ctx context.Context, reminder model.Reminder) (model.Reminder, error)
	UpdateReminder(ctx context.Context, reminder model.Reminder) error
	DeleteReminder(ctx context.Context, reminderID int64) error

	GetRemindersByHabitID(ctx context.Context, habitID int64) ([]model.Reminder, error)
	GetRemindersByUserID(ctx context.Context, userID int64) ([]model.Reminder, error)
}

type ProgressRepo interface {
	CreateProgress(ctx context.Context, progress model.Progress) (model.Progress, error)
	UpdateProgress(ctx context.Context, progress model.Progress) error
	DeleteProgress(ctx context.Context, habitID int64) error

	GetProgressByHabitID(ctx context.Context, habitID int64) (model.Progress, error)
	GetProgressByUserID(ctx context.Context, userID int64) ([]model.Progress, error)
}
