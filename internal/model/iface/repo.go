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
