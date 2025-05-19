package models

import (
    "time"

)


type CommentCreateRequestDTO struct {
    //UserID  uuid.UUID `json:"user_id" validate:"required"`
    VideoID uint `json:"video_id" validate:"required"`
    Content string    `json:"content" validate:"required"`
}

type CommentResponseDTO struct {
    ID       uint `json:"id"`
    UserID   uint `json:"user_id"`
    VideoID  uint `json:"video_id"`
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
