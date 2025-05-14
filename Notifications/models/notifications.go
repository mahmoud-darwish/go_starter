package models

import (
	"time"


)


type Notification struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"not null" json:"user_id"`
	ChannelID *uint      `json:"channel_id,omitempty"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	Source    string     `gorm:"type:varchar(255);not null" json:"source"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
