package entity

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
	CreatedAt  time.Time
}
