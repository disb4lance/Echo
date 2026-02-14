package entity

import (
	"time"

	"github.com/google/uuid"
)

type Boost struct {
	UserID     uuid.UUID
	Multiplier float64
	StartsAt   time.Time
	EndsAt     time.Time
	CreatedAt  time.Time
}
