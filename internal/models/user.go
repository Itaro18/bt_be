package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
    UserID   string `gorm:"primaryKey"`
    Name     string `json:"name"`
    Phone    string `gorm:"unique" json:"phone"` // unique if phone must be one per user
    Password string `json:"-"`
}

func (p *User) BeforeCreate(tx *gorm.DB) (err error) {
    p.UserID = uuid.NewString()
    return
}