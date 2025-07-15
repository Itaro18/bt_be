package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
    CustomerID string `gorm:"primaryKey"` // custom string ID
    Name       string
    Phone      string
}

func (p *Customer) BeforeCreate(tx *gorm.DB) (err error) {
    p.CustomerID = uuid.NewString()
    return
}
