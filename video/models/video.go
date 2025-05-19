package models

import (
	"time"
)

type Video struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ChannelID   uint      `gorm:"not null" json:"channel_id"`
	Type        int       `gorm:"not null" json:"type"`
	Path        string    `json:"path"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"not null" json:"description"`
	TimeStamp   time.Time `gorm:"not null" json:"time_stamp"`
	Genre       string    `gorm:"not null" json:"genre"`
}

func (Video) TableName() string {
	return "video"
}

type VideoCreateRequestDTO struct {
	ChannelID   uint   `json:"channel_id" validate:"required"`
	Type        int    `json:"type" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Genre       string `json:"genre" validate:"required"`
}

type VideoUpdateRequestDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	Type        int    `json:"type"`
}

type VideoResponseDTO struct {
	ID          uint      `json:"id"`
	ChannelID   uint      `json:"channel_id"`
	Type        int       `json:"type"`
	Path        string    `json:"path"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TimeStamp   time.Time `json:"time_stamp"`
	Genre       string    `json:"genre"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type VideoListResponseDTO struct {
	Videos []VideoResponseDTO `json:"videos"`
	Total  int64              `json:"total"`
}

type VideoPaginationParams struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}
