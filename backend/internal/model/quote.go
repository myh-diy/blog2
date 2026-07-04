package model

import "time"

type Quote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
