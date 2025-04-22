package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationCreateRequestDTO struct {
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	ChannelID *uuid.UUID `json:"channel_id,omitempty"`
	Content   string     `json:"content" validate:"required"`
	Source    string     `json:"source" validate:"required"`
}

type NotificationUpdateRequestDTO struct {
	Content string `json:"content" validate:"required"`
	Source  string `json:"source" validate:"required"`
}

type NotificationResponseDTO struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	ChannelID *uuid.UUID `json:"channel_id,omitempty"`
	Content   string     `json:"content"`
	Source    string     `json:"source"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NotificationResponseDTOFromModel(n Notification) NotificationResponseDTO {
	return NotificationResponseDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		ChannelID: n.ChannelID,
		Content:   n.Content,
		Source:    n.Source,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
