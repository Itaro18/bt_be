package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Property struct{
	PropertyID  string `gorm:"primaryKey"` 
	Name 		string
	City		string
}

func (p *Property) BeforeCreate(tx *gorm.DB) (err error) {
    p.PropertyID = uuid.NewString()
    return
}