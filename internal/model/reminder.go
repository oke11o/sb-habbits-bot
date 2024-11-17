package model

type Reminder struct {
	ID      int64       `db:"id" json:"id"`
	HabitID int64       `db:"habit_id" json:"habit_id"`
	UserID  int64       `db:"user_id" json:"user_id"`
	Time    string      `db:"time" json:"time"` // Время напоминания (формат HH:MM)
	Days    StringSlice `db:"days" json:"days"` // Дни недели (например, "mon,tue,wed")
}
