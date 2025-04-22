package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	ChannelID *uuid.UUID `gorm:"type:uuid;not null" json:"channel_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
