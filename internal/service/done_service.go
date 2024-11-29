package service

import (
	"context"
	"github.com/oke11o/sb-habits-bot/internal/model/iface"
)

func NewDone(habitRepo iface.HabitRepo, recordRepo iface.RecordRepo) *Done {
	return &Done{
		parser: &RecordParser{
			habitRepo: habitRepo,
		},
		recordRepo: recordRepo,
	}
}

type Done struct {
	parser     *RecordParser
	recordRepo iface.RecordRepo
}

func (s *Done) Done(ctx context.Context, userID int64, msg string) error {
	record, err := s.parser.ParseCommand(ctx, userID, msg)
	if err != nil {
		return err //TODO:
	}
	record, err = s.recordRepo.CreateRecord(ctx, record)
	if err != nil {
		return err //TODO:
	}

	return nil
}
