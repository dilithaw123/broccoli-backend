package group

import (
	"context"
	"errors"
)

var (
	ErrGroupExists   = errors.New("group already exists")
	ErrGroupNotFound = errors.New("group not found")
)

type GroupService interface {
	CreateUpdateGroup(ctx context.Context, g Group, userEmail string) error
	GetGroup(ctx context.Context, id int) (Group, error)
	GetGroupByName(ctx context.Context, name string) (Group, error)
	GetGroupsByEmail(ctx context.Context, email string) ([]Group, error)
}
