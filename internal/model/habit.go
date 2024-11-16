package model

import "time"

type Habit struct {
	ID         int64     `db:"id" json:"id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	Name       string    `db:"name" json:"name"`
	Type       string    `db:"type" json:"type"`
	Target     int64     `db:"target" json:"target"`
	TargetTime string    `db:"target_time" json:"target_time"`
	MaxTime    string    `db:"max_time" json:"max_time"`
	Unit       string    `db:"unit" json:"unit"`
	Points     int64     `db:"points" json:"points"`
	PointsMode string    `db:"points_mode" json:"points_mode"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
