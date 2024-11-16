package model

import "time"

type Record struct {
	ID        int64     `db:"id" json:"id"`
	HabitID   int64     `db:"habit_id" json:"habit_id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Value     int64     `db:"value" json:"value"` // Значение выполнения (количество повторений, длительность и т.д.)
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Points    int64     `db:"points" json:"points"` // Начисленные баллы
}
