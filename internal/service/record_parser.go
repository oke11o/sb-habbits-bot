package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/oke11o/sb-habits-bot/internal/model"
	"github.com/oke11o/sb-habits-bot/internal/model/iface"
)

var (
	ErrHabitNotFound  = errors.New("habit not found")
	ErrInvalidCommand = errors.New("invalid command format")
	ErrInvalidValue   = errors.New("invalid value format")
	ErrInvalidTime    = errors.New("invalid time format")
)

// ParseCommand принимает строку команды и возвращает Record для сохранения.
func ParseCommand(ctx context.Context, habitRepo iface.HabitRepo, userID int64, command string) (model.Record, error) {
	// Удаляем лишние пробелы и разбиваем команду на части
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return model.Record{}, ErrInvalidCommand
	}

	// Проверяем, что команда начинается с "done"
	if parts[0] != "done" {
		return model.Record{}, ErrInvalidCommand
	}

	// Извлекаем название привычки
	habitName := strings.Join(parts[1:len(parts)-1], " ")
	valueStr := parts[len(parts)-1]

	// Ищем привычку в базе данных
	habit, err := habitRepo.GetHabitByName(ctx, userID, habitName)
	if err != nil {
		return model.Record{}, fmt.Errorf("habitRepo.GetHabitByName() err: %w", ErrHabitNotFound)
	}

	// Создаём объект Record
	record := model.Record{
		HabitID:   habit.ID,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	// Обрабатываем значение в команде
	switch habit.Type {
	case "simple":
		// Простая привычка (без значений)
		record.Value = 1
		record.Points = habit.Points

	case "counter", "cumulative":
		// Привычки с количеством (отжимания, шаги)
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return model.Record{}, ErrInvalidValue
		}
		record.Value = value
		record.Points = value * habit.Points / habit.Target

	case "time":
		// Привычка по времени (проснуться в 6 утра)
		execTime, err := time.Parse("15:04", valueStr)
		if err != nil {
			return model.Record{}, ErrInvalidTime
		}
		targetTime, _ := time.Parse("15:04", habit.TargetTime)
		maxTime, _ := time.Parse("15:04", habit.MaxTime)

		// Расчёт баллов на основе времени
		if execTime.Before(targetTime) || execTime.Equal(targetTime) {
			record.Points = habit.Points
		} else if execTime.After(maxTime) {
			record.Points = 0
		} else {
			diff := maxTime.Sub(targetTime).Minutes()
			points := habit.Points * int64(maxTime.Sub(execTime).Minutes()) / int64(diff)
			record.Points = points
		}
		record.Value = 1

	default:
		return model.Record{}, fmt.Errorf("unsupported habit type: %s", habit.Type)
	}

	return record, nil
}
