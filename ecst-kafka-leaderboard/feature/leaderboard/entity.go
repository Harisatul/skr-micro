package leaderboard

import (
	"github.com/google/uuid"
)

type Grade struct {
	ID        uuid.UUID
	CreatedAt int64
	UpdatedAt int64
	Token     string
	TryoutID  uuid.UUID
	Version   int64
	Score     int64
}
