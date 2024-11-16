package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/oke11o/sb-habits-bot/internal/config"
	"github.com/oke11o/sb-habits-bot/internal/model"
)

func NewSessionRepo(cfg config.SqliteConfig) (*SessionRepo, error) {
	dbx, err := NewDb(cfg)
	if err != nil {
		return nil, err
	}

	return &SessionRepo{db: dbx}, nil
}

func NewSessionRepoWithDB(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

type SessionRepo struct {
	db *sqlx.DB
}

func (r *SessionRepo) SaveSession(ctx context.Context, session model.Session) (model.Session, error) {
	err := session.PrepareToSave()
	if err != nil {
		return session, fmt.Errorf("session.PrepareToSave() err: %w", err)
	}
	q := `insert into session (user_id,data,status,created_at,updated_at)
values (:user_id,:data,:status,:created_at,:updated_at)`
	raw, err := r.db.NamedExecContext(ctx, q, session)
	if err != nil {
		return session, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := raw.LastInsertId()
	if err != nil {
		return session, fmt.Errorf("raw.LastInsertId() err: %w", err)
	}
	session.ID = id

	return session, nil
}

func (r *SessionRepo) CloseSession(ctx context.Context, session model.Session) error {
	q := `update session set closed=1,updated_at=? where user_id=?`
	_, err := r.db.ExecContext(ctx, q, time.Now().Format(time.RFC3339), session.UserID)
	if err != nil {
		return fmt.Errorf("db.ExecContext() err: %w", err)
	}
	return nil
}

func (r *SessionRepo) GetOpenedSession(ctx context.Context, userID int64) (model.Session, error) {
	ses := model.Session{}
	q := `select id,user_id,data,status,created_at,updated_at from session where user_id=? and closed=0 order by id desc limit 1`
	err := r.db.GetContext(ctx, &ses, q, userID)
	if err == nil {
		err = ses.AfterGet()
		if err != nil {
			return ses, fmt.Errorf("session.AfterGet() err: %w", err)
		}
		return ses, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ses, nil
	}

	return ses, fmt.Errorf("db.GetContext() err: %w", err)
}
