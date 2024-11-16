package sqlite

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/oke11o/sb-habits-bot/internal/model"
)

type RecordRepo struct {
	db *sqlx.DB
}

func NewRecordRepo(db *sqlx.DB) *RecordRepo {
	return &RecordRepo{db: db}
}

func (r *RecordRepo) CreateRecord(ctx context.Context, record model.Record) (model.Record, error) {
	query := `
		INSERT INTO records (habit_id, user_id, value, timestamp, points)
		VALUES (:habit_id, :user_id, :value, :timestamp, :points)
		RETURNING id
	`

	res, err := r.db.NamedExecContext(ctx, query, record)
	if err != nil {
		return model.Record{}, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Record{}, fmt.Errorf("res.LastInsertId() err: %w", err)
	}
	record.ID = id
	return record, nil
}

func (r *RecordRepo) DeleteRecord(ctx context.Context, recordID int64) error {
	query := `DELETE FROM records WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, recordID)
	if err != nil {
		return fmt.Errorf("db.ExecContext() err: %w", err)
	}
	return nil
}

func (r *RecordRepo) GetRecordsByHabitID(ctx context.Context, habitID int64, limit int) ([]model.Record, error) {
	query := `SELECT * FROM records WHERE habit_id = $1 ORDER BY timestamp DESC LIMIT $2`
	var records []model.Record
	err := r.db.SelectContext(ctx, &records, query, habitID, limit)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext() err: %w", err)
	}
	return records, nil
}

func (r *RecordRepo) GetRecordsByUserID(ctx context.Context, userID int64, limit int) ([]model.Record, error) {
	query := `SELECT * FROM records WHERE user_id = $1 ORDER BY timestamp DESC LIMIT $2`
	var records []model.Record
	err := r.db.SelectContext(ctx, &records, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext() err: %w", err)
	}
	return records, nil
}

func (r *RecordRepo) GetLatestRecordByHabitID(ctx context.Context, habitID int64) (model.Record, error) {
	query := `SELECT * FROM records WHERE habit_id = $1 ORDER BY timestamp DESC LIMIT 1`
	var record model.Record
	err := r.db.GetContext(ctx, &record, query, habitID)
	if err != nil {
		return record, fmt.Errorf("db.GetContext() err: %w", err)
	}
	return record, nil
}
