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
	transaction, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer transaction.Rollback(ctx)
	// Session is new, create
	if err = pgxscan.Get(
		ctx,
		transaction,
		&id,
		"INSERT INTO sessions (group_id, create_date, shuffle_seed) VALUES ($1, $2, $3) RETURNING id",
		s.GroupID,
		s.CreateDate,
		s.ShuffleSeed,
	); err != nil {
		return 0, err
	}

	if _, err = transaction.Exec(
		ctx,
		`
		WITH prev_session AS (
			SELECT id 
			FROM sessions 
			WHERE id != $1
			AND group_id = $2
			ORDER BY id DESC 
			LIMIT 1
		)
		INSERT INTO user_submissions (user_id, session_id, yesterday, today, blockers)
		SELECT user_id, $1, today, '{}', '{}'
		FROM user_submissions
		WHERE session_id = (SELECT id FROM prev_session)
		AND EXISTS (SELECT 1 FROM prev_session);
		`,
		id,
		s.GroupID,
	); err != nil {
		return id, err
	}

	err = transaction.Commit(ctx)
	return id, err
}

func (repo *PgSessionRepo) UpdateShuffle(ctx context.Context, id uint64, seed uint16) error {
	conn, err := repo.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	_, err = conn.Exec(
		ctx,
		"UPDATE sessions SET shuffle_seed = $2 WHERE ID = $1",
		id,
		seed,
	)
	return err
}
