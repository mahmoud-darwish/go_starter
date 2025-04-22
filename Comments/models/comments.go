package models

import (
    "time"

    "github.com/google/uuid"
)

type Comment struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    VideoID   uuid.UUID `gorm:"type:uuid;not null" json:"video_id"`
    Content   string    `gorm:"type:text;not null" json:"content"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
