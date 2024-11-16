package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/oke11o/sb-habits-bot/internal/model"
)

type ProgressRepo struct {
	db *sqlx.DB
}

func NewProgressRepo(db *sqlx.DB) *ProgressRepo {
	return &ProgressRepo{db: db}
}

func (r *ProgressRepo) CreateProgress(ctx context.Context, progress model.Progress) (model.Progress, error) {
	query := `
		INSERT INTO progress (habit_id, user_id, accumulated_value, target, last_updated)
		VALUES (:habit_id, :user_id, :accumulated_value, :target, :last_updated)
	`
	res, err := r.db.NamedExecContext(ctx, query, progress)
	if err != nil {
		return model.Progress{}, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Progress{}, fmt.Errorf("res.LastInsertId() err: %w", err)
	}
	progress.ID = id
	return progress, nil
}

func (r *ProgressRepo) UpdateProgress(ctx context.Context, progress model.Progress) error {
	query := `
		UPDATE progress
		SET accumulated_value = :accumulated_value, target = :target, last_updated = :last_updated
		WHERE habit_id = :habit_id
	`

	_, err := r.db.NamedExecContext(ctx, query, progress)
	if err != nil {
		return fmt.Errorf("db.QueryRowxContext() err: %w", err)
	}
	return nil
}

func (r *ProgressRepo) DeleteProgress(ctx context.Context, habitID int64) error {
	query := `DELETE FROM progress WHERE habit_id = $1`
	_, err := r.db.ExecContext(ctx, query, habitID)
	if err != nil {
		return fmt.Errorf("db.ExecContext() err: %w", err)
	}
	return nil
}

func (r *ProgressRepo) GetProgressByHabitID(ctx context.Context, habitID int64) (model.Progress, error) {
	query := `SELECT * FROM progress WHERE habit_id = $1`
	var progress model.Progress
	err := r.db.GetContext(ctx, &progress, query, habitID)
	if err != nil {
		return progress, fmt.Errorf("db.GetContext() err: %w", err)
	}
	return progress, nil
}

func (r *ProgressRepo) GetProgressByUserID(ctx context.Context, userID int64) ([]model.Progress, error) {
	query := `SELECT * FROM progress WHERE user_id = $1`
	var progressList []model.Progress
	err := r.db.SelectContext(ctx, &progressList, query, userID)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext() err: %w", err)
	}
	return progressList, nil
}
