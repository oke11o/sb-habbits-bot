package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Habit struct {
	ID             int64       `db:"id" json:"id"`
	UserID         int64       `db:"user_id" json:"user_id"`
	Name           string      `db:"name" json:"name"`
	Type           string      `db:"type" json:"type"`
	Target         int64       `db:"target" json:"target"`
	TargetTime     string      `db:"target_time" json:"target_time"`
	MaxTime        string      `db:"max_time" json:"max_time"`
	Unit           string      `db:"unit" json:"unit"`
	Points         int64       `db:"points" json:"points"`
	PointsMode     string      `db:"points_mode" json:"points_mode"`
	CreatedAt      time.Time   `db:"created_at" json:"created_at"`
	TargetDuration string      `db:"target_duration" json:"target_duration"`
	IntervalDays   int64       `db:"interval_days" json:"interval_days"`
	Tasks          StringSlice `db:"tasks" json:"tasks"`
	Options        StringSlice `db:"options" json:"options"`
}

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal() err: %w", err)
	}
	return string(data), nil
}

func (s *StringSlice) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, s); err != nil {
			return fmt.Errorf("json.Unmarshal() err: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(v), s); err != nil {
			return fmt.Errorf("json.Unmarshal() err: %w", err)
		}
	default:
		return fmt.Errorf("failed to convert value to StringSlice: unsupported type %T", value)
	}
	return nil
}
