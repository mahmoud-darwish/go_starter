package models

import (
	"time"


)


type Channel struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"not null" json:"user_id"`
	Name   string     `gorm:"type:varchar(255);not null" json:"name"`
	Logo   string     `gorm:"type:varchar(255);not null" json:"logo"`
	Bio   string     `gorm:"type:text;not null" json:"bio"`
	TimeStamp time.Time `gorm:"type:timestamptz;not null" json:"time_stamp"`
}

func (Channel) TableName() string {
    return "channel"
}