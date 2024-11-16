package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type SessionStatus string

const (
	SessionStatusClosed   SessionStatus = "closed"
	SessionMembersProcess SessionStatus = "members_process"
)

type Session struct {
	ID        int64             `db:"id"`
	UserID    int64             `db:"user_id"`
	Data      string            `db:"data"`
	Closed    bool              `db:"closed"`
	dataArgs  map[string]string `db:"-"`
	Status    SessionStatus     `db:"status"`
	CreatedAt string            `db:"created_at"`
	UpdatedAt string            `db:"updated_at"`
}

func (s *Session) GetArg(key string) (string, bool) {
	if s.dataArgs == nil {
		return "", false
	}
	value, ok := s.dataArgs[key]
	return value, ok
}

func (s *Session) SetArg(key string, value string) {
	if s.dataArgs == nil {
		s.dataArgs = make(map[string]string)
	}
	s.dataArgs[key] = value
}

func (s *Session) RemoveArg(key string) {
	delete(s.dataArgs, key)
}

func (s *Session) SetStatus(status SessionStatus) {
	s.Status = status
}

func (s *Session) PrepareToSave() error {
	if s.dataArgs == nil {
		s.Data = "{}"
		return nil
	}
	b, err := json.Marshal(s.dataArgs)
	if err != nil {
		return fmt.Errorf("json.Marshal() err: %w", err)
	}
	s.Data = string(b)
	return nil
}

func (s *Session) AfterGet() error {
	var data map[string]string
	err := json.Unmarshal([]byte(s.Data), &data)
	if err != nil {
		return fmt.Errorf("json.Unmarshal() err: %w", err)
	}
	s.dataArgs = data
	return nil
}

func NewMembersSession(userID int64) Session {
	return newSession(userID, SessionMembersProcess)
}

func newSession(userID int64, process SessionStatus) Session {
	return Session{
		UserID:    userID,
		Status:    process,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}
