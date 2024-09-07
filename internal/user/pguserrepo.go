package user

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const UNIQUE_VIOLATION = "23505"

type PgUserRepo struct {
	db *pgxpool.Pool
}

func NewPgUserRepo(db *pgxpool.Pool) *PgUserRepo {
	return &PgUserRepo{db: db}
}

func (repo *PgUserRepo) CreateUser(ctx context.Context, u User) (User, error) {
	var id uint64
	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		return User{}, err
	}
	defer conn.Release()
	err = pgxscan.Get(
		ctx,
		conn,
		&id,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		u.Name,
		u.Email,
	)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_VIOLATION {
			return User{}, ErrUserAlreadyExists
		}
		return User{}, err
	}
	u.ID = id
	return u, nil
}

func (repo *PgUserRepo) GetUserByID(ctx context.Context, id uint64) (User, error) {
	var u User
	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		return User{}, err
	}
	defer conn.Release()
	err = pgxscan.Get(
		ctx,
		conn,
		&u,
		"SELECT id, name, email FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

func (repo *PgUserRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		return User{}, err
	}
	err = pgxscan.Get(
		ctx,
		conn,
		&u,
		"SELECT id, name, email FROM users WHERE email = $1",
		email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

func (repo *PgUserRepo) GetUserSubmission(
	ctx context.Context,
	userId, sessionId uint64,
) (UserSubmission, error) {
	var us UserSubmission
	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		return UserSubmission{}, err
	}
	defer conn.Release()
	err = pgxscan.Get(
		ctx,
		conn,
		&us,
		`SELECT * FROM user_submissions WHERE user_id = $1 AND session_id = $2`,
		userId,
		sessionId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserSubmission{}, ErrorUserSubmissionNotFound
		}
		return UserSubmission{}, err
	}

	return us, nil
}

func (repo *PgUserRepo) GetAllUserSubmissionsForSession(
	ctx context.Context,
	sessionId uint64,
) ([]DBUserSubmission, error) {
	var us []DBUserSubmission
	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	err = pgxscan.Select(
		ctx,
		conn,
		&us,
		`SELECT us.*, u.name FROM user_submissions us JOIN users u ON us.user_id = u.id WHERE us.session_id = $1`,
		sessionId,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return us, nil
		}
		return nil, err
	}
	return us, nil
}

func (repo *PgUserRepo) CreateUpdateUserSubmission(ctx context.Context, us UserSubmission) error {
	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(
		ctx,
		`MERGE INTO user_submissions us
		USING (SELECT $1::BIGINT as user_id, $2::BIGINT as session_id) s
		ON us.user_id = s.user_id AND us.session_id = s.session_id
		WHEN MATCHED THEN
			UPDATE SET yesterday = $3, today = $4, blockers = $5
		WHEN NOT MATCHED THEN
			INSERT (user_id, session_id, yesterday, today, blockers) VALUES (s.user_id, s.session_id, $3, $4, $5)`,
		us.UserId,
		us.SessionId,
		us.Yesterday,
		us.Today,
		us.Blockers,
	)
	return err
}
