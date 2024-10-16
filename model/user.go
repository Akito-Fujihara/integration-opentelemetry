package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `json:"name" gorm:"size:100"`
	Email     string `json:"email" gorm:"unique;size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
