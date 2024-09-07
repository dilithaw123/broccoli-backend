package user

import (
	"context"
	"errors"
)

var (
	ErrUserAlreadyExists        = errors.New("user with email already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrorUserSubmissionNotFound = errors.New("user submission not found")
)

type UserService interface {
	CreateUser(ctx context.Context, u User) (User, error)
	GetUserByID(ctx context.Context, id uint64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserSubmission(ctx context.Context, userId, sessionId uint64) (UserSubmission, error)
	GetAllUserSubmissionsForSession(
		ctx context.Context,
		sessionId uint64,
	) ([]DBUserSubmission, error)
	CreateUpdateUserSubmission(ctx context.Context, us UserSubmission) error
}
