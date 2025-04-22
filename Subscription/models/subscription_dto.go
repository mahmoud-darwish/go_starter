package models

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionCreateRequestDTO struct {
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	ChannelID *uuid.UUID `json:"channel_id,omitempty"`
}



type SubscriptionResponseDTO struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	ChannelID *uuid.UUID `json:"channel_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func SubscriptionResponseDTOFromModel(n Subscription) SubscriptionResponseDTO {
	return SubscriptionResponseDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		ChannelID: n.ChannelID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
