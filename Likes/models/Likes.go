package models

import (
	"time"

)

type Like struct {
    ID      uint `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID  uint `gorm:"not null" json:"user_id"`
    VideoID uint `gorm:"not null" json:"video_id"`    
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
