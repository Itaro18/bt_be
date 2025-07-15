package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	BookingID string `gorm:"primaryKey"`

	CustomerID string
	Customer   *Customer `gorm:"foreignKey:CustomerID;references:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	PhoneNo    string
	PropertyID string
	Property   *Property `gorm:"foreignKey:PropertyID;references:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Handler string

	CheckInDate  time.Time
	CheckOutDate time.Time
	TotalAmount  float64
	AdvancePaid  float64

	Status  string
	Floor   string
	Through string
	City    string
	Remarks string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	b.BookingID = uuid.NewString()
	return
}
