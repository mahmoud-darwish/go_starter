package models

import (
	"time"

	"github.com/google/uuid"
)

type LikeCreateRequestDTO struct {
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	VideoID uuid.UUID `json:"video_id" validate:"required"`
}



type LikeResponseDTO struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	VideoID uuid.UUID `json:"video_id" validate:"required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func LikeResponseDTOFromModel(n Like) LikeResponseDTO {
	return LikeResponseDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		VideoID: n.VideoID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
