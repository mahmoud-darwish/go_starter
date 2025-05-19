package models

import (
	"time"
)

type ChannelCreateRequestDTO struct {
	//UserID    uint   `json:"user_id" validate:"required"`
	Name   string `json:"Name" validate:"required"`
	Logo    string `json:"Logo" validate:"required"`
	Bio    string `json:"Bio" validate:"required"`
}

type ChannelUpdateRequestDTO struct {
	Name string `json:"name" validate:"required"`
	Logo string `json:"Logo" validate:"required"`
	Bio  string `json:"Bio" validate:"required"`
}

type ChannelResponseDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Name   string    `json:"Name"`
	Logo    string    `json:"Logo"`
	Bio    string    `json:"Bio"`
	TimeStamp time.Time `json:"TimeStamp"`

}

func ChannelResponseDTOFromModel(n Channel) ChannelResponseDTO {
	return ChannelResponseDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		Name: n.Name,
		Logo:   n.Logo,
		Bio:    n.Bio,
		TimeStamp: n.TimeStamp,
	}
}
