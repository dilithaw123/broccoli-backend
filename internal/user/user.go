package user

import (
	"encoding/json"
	"strings"
)

type User struct {
	ID    uint64 `json:"id"    db:"id"`
	Name  string `json:"name"  db:"name"`
	Email string `json:"email" db:"email"`
}

type UserSubmission struct {
	ID        uint64   `json:"id"         db:"id"`
	UserId    uint64   `json:"user_id"    db:"user_id"`
	SessionId uint64   `json:"session_id" db:"session_id"`
	Yesterday []string `json:"yesterday"  db:"yesterday"`
	Today     []string `json:"today"      db:"today"`
	Blockers  []string `json:"blockers"   db:"blockers"`
}

type DBUserSubmission struct {
	UserSubmission
	Name string `json:"name" db:"name"`
}

func NewUser(name, email string) User {
	return User{
		Name:  name,
		Email: strings.ToLower(email),
	}
}

func NewUserSubmission(
	userId, sessionId uint64,
	name string,
	yesterday, today, blockers []string,
) UserSubmission {
	return UserSubmission{
		UserId:    userId,
		SessionId: sessionId,
		Yesterday: yesterday,
		Today:     today,
		Blockers:  blockers,
	}
}

func (u User) JSON() ([]byte, error) {
	return json.Marshal(u)
}
