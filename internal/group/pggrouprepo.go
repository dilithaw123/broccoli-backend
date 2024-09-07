package group

import (
	"context"
	"errors"

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

func (repo *PgGroupRepo) GetGroup(ctx context.Context, id int) (Group, error) {
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
	conn, err := repo.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return g, err
	}
	err = pgxscan.Get(
		ctx,
		conn,
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
