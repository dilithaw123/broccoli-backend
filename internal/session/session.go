package session

import (
	"time"

	"github.com/dilithaw123/broccoli-backend/internal/types"
)

type Session struct {
	ID         uint64           `json:"id"          db:"id"`
	GroupID    uint64           `json:"group_id"    db:"group_id"`
	CreateDate types.CustomTime `json:"create_date" db:"create_date"`
}

func NewSession(groupID uint64) Session {
	return Session{
		GroupID:    groupID,
		CreateDate: types.CustomTime(time.Now()),
	}
}
