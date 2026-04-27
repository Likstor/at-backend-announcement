package domain

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID          uint64
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedBy   uuid.UUID
	UpdatedAt   time.Time
	Title       string
	Description string
}
