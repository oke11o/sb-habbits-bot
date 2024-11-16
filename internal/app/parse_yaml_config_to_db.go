package app

import (
	"context"
	"fmt"
	"github.com/oke11o/sb-habits-bot/internal/model/iface"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"time"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/model"
	"github.com/oke11o/sb-habits-bot/internal/repository/sqlite"
)

func ParseYAMLConfigToDB(ctx context.Context, cfg config.Config, userID int64, filePath string, l *slog.Logger) error {
	db, err := sqlite.NewDb(cfg.Sqlite)
	if err != nil {
		return fmt.Errorf("sqlite.NewDb: %w", err)
	}
	defer db.Close()

	habitCfg, err := parseYAML(filePath)
	if err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	habitRepo := sqlite.NewHabitRepoWithDB(db)
	reminderRepo := sqlite.NewReminderRepo(db)
	err = addHabitsToDB(ctx, habitRepo, reminderRepo, userID, habitCfg, l)
	if err != nil {
		return fmt.Errorf("add habits to db: %w", err)
	}

	return nil
}

func addHabitsToDB(ctx context.Context, habitRepo iface.HabitRepo, reminderRepo iface.ReminderRepo, userID int64, config Config, l *slog.Logger) error {
	for _, habit := range config.Habits {
		// Создаём запись привычки
		habitRecord := model.Habit{
			UserID:     userID,
			Name:       habit.Name,
			Type:       habit.Type,
			Target:     habit.Target,
			TargetTime: habit.TargetTime,
			MaxTime:    habit.MaxTime,
			Points:     habit.Points,
			PointsMode: habit.PointsMode,
			Unit:       habit.Unit,
			CreatedAt:  time.Now(),
		}

		habitRecord, err := habitRepo.CreateHabit(ctx, habitRecord)
		if err != nil {
			return fmt.Errorf("habitRepo.CreateHabit() err: %w", err)
		}

		if habit.Reminder.Time != "" {
			reminder := model.Reminder{
				HabitID: habitRecord.ID,
				UserID:  userID,
				Time:    habit.Reminder.Time,
				Days:    fmt.Sprintf("%v", habit.Reminder.Days),
			}

			_, err := reminderRepo.CreateReminder(ctx, reminder)
			if err != nil {
				return fmt.Errorf("reminderRepo.CreateReminder() err: %w", err)
			}
		}

		l.InfoContext(ctx, "Привычка успешно добавлена", slog.String("habit", habit.Name), slog.Int64("habit_id", habitRecord.ID))
	}

	return nil
}

func parseYAML(filePath string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("os.ReadFile() err: %w", err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("yaml.Unmarshal() err: %w", err)
	}

	return config, nil
}

type HabitConfig struct {
	Name           string   `yaml:"name"`
	Type           string   `yaml:"type"`
	Target         int64    `yaml:"target"`
	TargetTime     string   `yaml:"target_time"`
	MaxTime        string   `yaml:"max_time"`
	TargetDuration string   `yaml:"target_duration"`
	Points         int64    `yaml:"points"`
	PointsMode     string   `yaml:"points_mode"`
	Unit           string   `yaml:"unit"`
	IntervalDays   int64    `yaml:"interval_days"`
	Tasks          []string `yaml:"tasks"`
	Options        []string `yaml:"options"`
	Reminder       struct {
		Time string   `yaml:"time"`
		Days []string `yaml:"days"`
	} `yaml:"reminder"`
}

type Config struct {
	Habits []HabitConfig `yaml:"habits"`
}
