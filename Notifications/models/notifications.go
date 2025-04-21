package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ChannelID *uuid.UUID `gorm:"type:uuid" json:"channel_id,omitempty"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	Source    string     `gorm:"type:varchar(255);not null" json:"source"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
