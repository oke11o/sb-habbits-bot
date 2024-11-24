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

var now = time.Now

var (
	ErrHabitNotFound  = errors.New("habit not found")
	ErrInvalidCommand = errors.New("invalid command format")
	ErrInvalidValue   = errors.New("invalid value format")
	ErrInvalidTime    = errors.New("invalid time format")
)

func NewRecordParser(habitRepo iface.HabitRepo) *RecordParser {
	return &RecordParser{
		habitRepo: habitRepo,
	}
}

type RecordParser struct {
	habitRepo iface.HabitRepo
}

func (p *RecordParser) ParseCommand(ctx context.Context, userID int64, command string) (model.Record, error) {
	// Удаляем лишние пробелы и разбиваем команду на части
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return model.Record{}, ErrInvalidCommand
	}

	if parts[0] != "/done" {
		return model.Record{}, ErrInvalidCommand
	}

	habitName := strings.Join(parts[1:len(parts)], " ")
	valueStr := ""

	habit, err := p.habitRepo.GetHabitByName(ctx, userID, habitName)
	if err != nil {
		habitName = strings.Join(parts[1:len(parts)-1], " ")
		valueStr = parts[len(parts)-1]
		habit, err = p.habitRepo.GetHabitByName(ctx, userID, habitName)
		if err != nil {
			return model.Record{}, fmt.Errorf("habitRepo.GetHabitByName() err: %w", ErrHabitNotFound)
		}
	}

	record := model.Record{
		HabitID:   habit.ID,
		UserID:    userID,
		Timestamp: now(),
	}

	switch habit.Type {
	case "simple":
		record.Value = 1
		record.Points = habit.Points

	case "counter", "cumulative":
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return model.Record{}, ErrInvalidValue
		}
		record.Value = value
		if value < habit.Target {
			record.Points = value * habit.Points / habit.Target
		} else {
			record.Points = habit.Points
		}

	case "time":
		execTime, err := time.Parse("15:04", valueStr)
		if err != nil {
			return model.Record{}, ErrInvalidTime
		}
		targetTime, _ := time.Parse("15:04", habit.TargetTime)
		maxTime, _ := time.Parse("15:04", habit.MaxTime)

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
