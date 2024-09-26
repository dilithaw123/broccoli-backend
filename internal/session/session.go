package session

import (
	"math/rand/v2"
	"time"

	"github.com/dilithaw123/broccoli-backend/internal/types"
)

type Session struct {
	ID          uint64           `json:"id"           db:"id"`
	GroupID     uint64           `json:"group_id"     db:"group_id"`
	CreateDate  types.CustomTime `json:"create_date"  db:"create_date"`
	ShuffleSeed uint16           `json:"shuffle_seed" db:"shuffle_seed"`
}

func NewSession(groupID uint64) Session {
	return Session{
		GroupID:     groupID,
		CreateDate:  types.CustomTime(time.Now()),
		ShuffleSeed: NewSeed(),
	}
}

func NewSeed() uint16 {
	return uint16(rand.UintN(32767))
}
