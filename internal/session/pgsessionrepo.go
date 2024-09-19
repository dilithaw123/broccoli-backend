package session

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const UNIQUE_VIOLATION = "23505"

type PgSessionRepo struct {
	db *pgxpool.Pool
}

func NewPgSessionRepo(db *pgxpool.Pool) *PgSessionRepo {
	return &PgSessionRepo{db: db}
}

func (repo *PgSessionRepo) GetSession(ctx context.Context, id uint64) (Session, error) {
	var s Session
	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		return s, err
	}
	defer conn.Release()
	err = pgxscan.Get(
		ctx,
		conn,
		&s,
		"SELECT * FROM sessions WHERE id = $1",
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Session{}, ErrSessionNotFound
		}
		return Session{}, err
	}
	return s, nil
}

// Get most recent session by group ID
func (repo *PgSessionRepo) GetSessionByGroupID(
	ctx context.Context,
	groupID uint64,
) (Session, error) {
	var s Session
	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		return s, err
	}
	defer conn.Release()
	err = pgxscan.Get(
		ctx,
		conn,
		&s,
		"SELECT * FROM sessions WHERE group_id = $1 ORDER BY create_date DESC LIMIT 1",
		groupID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Session{}, ErrSessionNotFound
		}
		return Session{}, err
	}
	return s, nil
}

func (repo *PgSessionRepo) CreateSession(ctx context.Context, s Session) (uint64, error) {
	conn, err := repo.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return 0, err
	}
	var id uint64
	err = pgxscan.Get(
		ctx,
		conn,
		&id,
		`SELECT s.id from sessions s
		JOIN groups g ON s.group_id = g.id
		WHERE s.group_id = $1
		AND (s.create_date at time zone g.timezone)::date = ($2 at time zone g.timezone)::date`,
		s.GroupID,
		s.CreateDate,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	if id != 0 {
		return id, nil
	}
	err = pgxscan.Get(
		ctx,
		conn,
		&id,
		"INSERT INTO sessions (group_id, create_date) VALUES ($1, $2) RETURNING id",
		s.GroupID,
		s.CreateDate,
	)
	return id, err
}
