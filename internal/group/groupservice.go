package group

import (
	"context"
	"errors"
)

var (
	ErrGroupExists      = errors.New("group already exists")
	ErrGroupNotFound    = errors.New("group not found")
	ErrUserNotPermitted = errors.New("user does not have permission")
)

type GroupService interface {
	CreateUpdateGroup(ctx context.Context, g Group, userEmail string) error
	GetGroup(ctx context.Context, id uint64) (Group, error)
	GetGroupByName(ctx context.Context, name string) (Group, error)
	GetGroupsByEmail(ctx context.Context, email string) ([]Group, error)
	GroupContainsUser(ctx context.Context, groupID uint64, userEmail string) (bool, error)
	AddUserToGroup(ctx context.Context, groupID uint64, userEmail string) error
	DeleteGroup(ctx context.Context, id uint64, userEmail string) error
}
