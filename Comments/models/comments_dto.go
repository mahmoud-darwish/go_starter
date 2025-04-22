package models

import (
    "time"

    "github.com/google/uuid"
)


type CommentCreateRequestDTO struct {
    UserID  uuid.UUID `json:"user_id" validate:"required"`
    VideoID uuid.UUID `json:"video_id" validate:"required"`
    Content string    `json:"content" validate:"required"`
}

type CommentResponseDTO struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    VideoID   uuid.UUID `json:"video_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func CommentResponseDTOFromModel(comment Comment) CommentResponseDTO {
    return CommentResponseDTO{
        ID:        comment.ID,
        UserID:    comment.UserID,
        VideoID:   comment.VideoID,
        Content:   comment.Content,
        CreatedAt: comment.CreatedAt,
        UpdatedAt: comment.UpdatedAt,
    }
}
