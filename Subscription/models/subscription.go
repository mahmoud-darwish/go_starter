package models

import (
	"time"

)

type Subscription struct {
	ID        uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint  `gorm:"not null" json:"user_id"`
	ChannelID *uint `gorm:"" json:"channel_id"`	
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
