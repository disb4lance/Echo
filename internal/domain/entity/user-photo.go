package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserPhoto struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	URL       string
	Position  int
	CreatedAt time.Time
}
