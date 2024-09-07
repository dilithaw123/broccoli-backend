package session

import (
	"context"
	"errors"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExists   = errors.New("session already exists")
)

type SessionService interface {
	GetSession(ctx context.Context, id uint64) (Session, error)
	GetSessionByGroupID(ctx context.Context, groupID uint64) (Session, error)
	CreateSession(ctx context.Context, s Session) (uint64, error)
}
