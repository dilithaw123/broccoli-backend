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

// If email is not in allowed_emails, insert or update the group
func (repo *PgGroupRepo) CreateUpdateGroup(
	ctx context.Context,
	g Group,
	userEmail string,
) error {
	conn, err := repo.db.Acquire(ctx)
	defer conn.Release()
	userEmail = strings.ToLower(userEmail)
	if err != nil {
		return err
	}
	_, err = conn.Exec(
		ctx,
		`MERGE into groups g
		USING (SELECT $1::text as name) s
		ON g.name = s.name AND $2 = ANY(g.allowed_emails)
		WHEN MATCHED THEN
			UPDATE SET allowed_emails = $3
		WHEN NOT MATCHED THEN
			INSERT (name, allowed_emails) VALUES ($1, $2, $3)`,
		g.Name,
		userEmail,
		g.AllowedEmails,
	)
	if err != nil {
		return err
	}
	return nil
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
	email = strings.ToLower(email)
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
	userEmail = strings.ToLower(userEmail)
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
