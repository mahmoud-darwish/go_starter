package models

import (
	"time"

)

type LikeCreateRequestDTO struct {
    UserID  uint `json:"user_id" validate:"required"`
    VideoID uint `json:"video_id" validate:"required"`
}
type LikeResponseDTO struct {
    ID        uint  `json:"id"`
    UserID    uint  `json:"user_id"`
    VideoID   uint  `json:"video_id"`
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
