package models

import (
	"time"

)

type SubscriptionCreateRequestDTO struct {
    UserID    uint  `json:"user_id" validate:"required"`
    ChannelID *uint `json:"channel_id,omitempty"`
}
type SubscriptionResponseDTO struct {
    ID        uint  `json:"id"`
    UserID    uint  `json:"user_id"`
    ChannelID *uint `json:"channel_id,omitempty"`
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
