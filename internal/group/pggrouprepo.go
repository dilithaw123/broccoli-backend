package group

import (
	"context"
	"errors"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const UNIQUE_VIOLATION = "23505"

type PgGroupRepo struct {
	db *pgxpool.Pool
}

func NewPgGroupRepo(db *pgxpool.Pool) *PgGroupRepo {
	return &PgGroupRepo{db: db}
}

func (repo *PgGroupRepo) CreateGroup(ctx context.Context, g Group) error {
	conn, err := repo.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	var exists bool
	err = pgxscan.Get(
		ctx,
		conn,
		&exists,
		"SELECT true FROM groups WHERE name = $1",
		g.Name,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if exists {
		return ErrGroupExists
	}
	_, err = conn.Exec(
		ctx,
		"INSERT INTO groups (name, allowed_emails,timezone) VALUES ($1, $2, $3)",
		g.Name,
		g.AllowedEmails,
		g.Timezone,
	)
	return err
}

func (repo *PgGroupRepo) GetGroup(ctx context.Context, id uint64) (Group, error) {
	var g Group
	conn, err := repo.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return Group{}, err
	}
	err = pgxscan.Get(
		ctx,
		conn,
		&g,
		"SELECT * FROM groups WHERE id = $1",
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Group{}, ErrGroupNotFound
		}
		return Group{}, err
	}
	return g, nil
}

func (repo *PgGroupRepo) GetGroupByName(ctx context.Context, name string) (Group, error) {
	var g Group
	err := pgxscan.Get(
		ctx,
		repo.db,
		&g,
		"SELECT * FROM groups WHERE name = $1",
		name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Group{}, ErrGroupNotFound
		}
		return Group{}, err
	}
	return g, nil
}

func (repo *PgGroupRepo) GetGroupsByEmail(ctx context.Context, email string) ([]Group, error) {
	var groups []Group
	conn, err := repo.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return groups, err
	}
	err = pgxscan.Select(
		ctx,
		conn,
		&groups,
		"SELECT * FROM groups WHERE $1 = ANY(allowed_emails)",
		email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return groups, ErrGroupNotFound
		}
		return groups, nil
	}
	return groups, nil
}

func (repo *PgGroupRepo) GroupContainsUser(
	ctx context.Context,
	groupID uint64,
	userEmail string,
) (bool, error) {
	var valid bool
	err := pgxscan.Get(
		ctx,
		repo.db,
		&valid,
		"SELECT $1 = ANY(allowed_emails) FROM groups WHERE id = $2",
		userEmail,
		groupID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrGroupNotFound
		}
		return false, err
	}
	return valid, nil
}

func (repo *PgGroupRepo) AddUserToGroup(
	ctx context.Context,
	groupID uint64,
	userEmail string,
) error {
	userEmail = strings.ToLower(userEmail)
	_, err := repo.db.Exec(
		ctx,
		`UPDATE groups SET allowed_emails = ARRAY_APPEND(allowed_emails, $1) WHERE id = $2 AND $1 <> ALL(allowed_emails)`,
		userEmail,
		groupID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PgGroupRepo) DeleteGroup(ctx context.Context, id uint64, userEmail string) error {
	conn, err := repo.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	var isAllowed bool
	err = pgxscan.Get(
		ctx,
		conn,
		&isAllowed,
		"SELECT $1 = ANY(allowed_emails) FROM groups WHERE id = $2",
		userEmail,
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrGroupNotFound
		}
		return err
	}
	if !isAllowed {
		return ErrUserNotPermitted
	}
	_, err = conn.Exec(
		ctx,
		"DELETE FROM groups WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
