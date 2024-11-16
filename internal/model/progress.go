package model

import "time"

type Progress struct {
	ID               int64     `db:"id" json:"id"`
	HabitID          int64     `db:"habit_id" json:"habit_id"`
	UserID           int64     `db:"user_id" json:"user_id"`
	AccumulatedValue int64     `db:"accumulated_value" json:"accumulated_value"` // Накопленное значение
	Target           int64     `db:"target" json:"target"`                       // Целевая величина
	LastUpdated      time.Time `db:"last_updated" json:"last_updated"`
}
