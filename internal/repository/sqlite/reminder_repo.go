package sqlite

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/oke11o/sb-habits-bot/internal/model"
)

type ReminderRepo struct {
	db *sqlx.DB
}

func NewReminderRepo(db *sqlx.DB) *ReminderRepo {
	return &ReminderRepo{db: db}
}

func (r *ReminderRepo) CreateReminder(ctx context.Context, reminder model.Reminder) (model.Reminder, error) {
	query := `
		INSERT INTO reminders (habit_id, user_id, time, days)
		VALUES (:habit_id, :user_id, :time, :days)
		RETURNING id
	`
	res, err := r.db.NamedExecContext(ctx, query, reminder)
	if err != nil {
		return model.Reminder{}, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Reminder{}, fmt.Errorf("res.LastInsertId() err: %w", err)
	}
	reminder.ID = id
	return reminder, nil
}

func (r *ReminderRepo) UpdateReminder(ctx context.Context, reminder model.Reminder) error {
	query := `
		UPDATE reminders
		SET time = :time, days = :days
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, reminder)
	if err != nil {
		return fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	return nil
}

func (r *ReminderRepo) DeleteReminder(ctx context.Context, reminderID int64) error {
	query := `DELETE FROM reminders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, reminderID)
	if err != nil {
		return fmt.Errorf("db.ExecContext() err: %w", err)
	}
	return nil
}

func (r *ReminderRepo) GetRemindersByHabitID(ctx context.Context, habitID int64) ([]model.Reminder, error) {
	query := `SELECT * FROM reminders WHERE habit_id = $1`
	var reminders []model.Reminder
	err := r.db.SelectContext(ctx, &reminders, query, habitID)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext() err: %w", err)
	}
	return reminders, nil
}

func (r *ReminderRepo) GetRemindersByUserID(ctx context.Context, userID int64) ([]model.Reminder, error) {
	query := `SELECT * FROM reminders WHERE user_id = $1`
	var reminders []model.Reminder
	err := r.db.SelectContext(ctx, &reminders, query, userID)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext() err: %w", err)
	}
	return reminders, nil
}
