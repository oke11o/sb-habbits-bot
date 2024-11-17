package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/model"
)

func NewHabitRepo(cfg config.SqliteConfig) (*HabitRepo, error) {
	dbx, err := NewDb(cfg)
	if err != nil {
		return nil, err
	}

	return &HabitRepo{db: dbx}, nil
}

func NewHabitRepoWithDB(db *sqlx.DB) *HabitRepo {
	return &HabitRepo{db: db}
}

type HabitRepo struct {
	db *sqlx.DB
}

func (r *HabitRepo) CreateHabit(ctx context.Context, habit model.Habit) (model.Habit, error) {
	query := `
		INSERT INTO habits (
			user_id, name, type, target, target_time, max_time, unit, points, points_mode,
			target_duration, interval_days, tasks, options, created_at
		)
		VALUES (
			:user_id, :name, :type, :target, :target_time, :max_time, :unit, :points, :points_mode,
			:target_duration, :interval_days, :tasks, :options, :created_at
		)
	`
	res, err := r.db.NamedExecContext(ctx, query, habit)
	if err != nil {
		return model.Habit{}, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Habit{}, fmt.Errorf("res.LastInsertId() err: %w", err)
	}
	habit.ID = id
	return habit, nil
}

func (r *HabitRepo) UpdateHabit(ctx context.Context, habit model.Habit) error {
	query := `
		UPDATE habits
		SET 
			name = :name, type = :type, target = :target, target_time = :target_time, 
			max_time = :max_time, unit = :unit, points = :points, points_mode = :points_mode,
			target_duration = :target_duration, interval_days = :interval_days, 
			tasks = :tasks, options = :options
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, habit)
	if err != nil {
		return fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	return nil
}

func (r *HabitRepo) UpsertHabit(ctx context.Context, habit model.Habit) (model.Habit, error) {
	query := `
		INSERT INTO habits (
			user_id, name, type, target, target_time, max_time, unit, points, points_mode,
			target_duration, interval_days, tasks, options, created_at
		)
		VALUES (
			:user_id, :name, :type, :target, :target_time, :max_time, :unit, :points, :points_mode,
			:target_duration, :interval_days, :tasks, :options, :created_at
		)
		ON CONFLICT(user_id, name) DO UPDATE SET
			type = EXCLUDED.type,
			target = EXCLUDED.target,
			target_time = EXCLUDED.target_time,
			max_time = EXCLUDED.max_time,
			unit = EXCLUDED.unit,
			points = EXCLUDED.points,
			points_mode = EXCLUDED.points_mode,
			target_duration = EXCLUDED.target_duration,
			interval_days = EXCLUDED.interval_days,
			tasks = EXCLUDED.tasks,
			options = EXCLUDED.options
	`

	res, err := r.db.NamedExecContext(ctx, query, habit)
	if err != nil {
		return model.Habit{}, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Habit{}, fmt.Errorf("res.LastInsertId() err: %w", err)
	}
	habit.ID = id
	return habit, nil
}

func (r *HabitRepo) GetHabitByID(ctx context.Context, habitID int64) (model.Habit, error) {
	query := `SELECT * FROM habits WHERE id = $1`
	var habit model.Habit
	err := r.db.GetContext(ctx, &habit, query, habitID)
	if err != nil {
		return model.Habit{}, fmt.Errorf("db.GetContext() err: %w", err)
	}
	return habit, nil
}

func (r *HabitRepo) GetHabitsByUserID(ctx context.Context, userID int64) ([]model.Habit, error) {
	query := `SELECT * FROM habits WHERE user_id = $1`
	var habits []model.Habit
	err := r.db.SelectContext(ctx, &habits, query, userID)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext() err: %w", err)
	}
	return habits, nil
}

func (r *HabitRepo) GetHabitByName(ctx context.Context, userID int64, name string) (model.Habit, error) {
	query := `SELECT * FROM habits WHERE user_id = $1 AND name = $2`
	var habit model.Habit
	err := r.db.GetContext(ctx, &habit, query, userID, name)
	if err != nil {
		return model.Habit{}, fmt.Errorf("db.GetContext() err: %w", err)
	}
	return habit, nil
}

func (r *HabitRepo) DeleteHabitByID(ctx context.Context, habitID int64) error {
	query := `DELETE FROM habits WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, habitID)
	if err != nil {
		return fmt.Errorf("db.ExecContext() err: %w", err)
	}
	return nil
}

func (r *HabitRepo) DeleteHabitByName(ctx context.Context, userID int64, name string) error {
	query := `DELETE FROM habits WHERE user_id = $1 AND name = $2`
	_, err := r.db.ExecContext(ctx, query, userID, name)
	if err != nil {
		return fmt.Errorf("db.ExecContext() err: %w", err)
	}
	return nil
}
