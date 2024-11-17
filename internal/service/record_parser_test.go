package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/oke11o/sb-habits-bot/internal/model"
	"github.com/oke11o/sb-habits-bot/internal/model/iface"
)

func TestParseCommand(t *testing.T) {
	type args struct {
		ctx       context.Context
		habitRepo iface.HabitRepo
		userID    int64
		command   string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Record
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommand(tt.args.ctx, tt.args.habitRepo, tt.args.userID, tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}
