package entity

import "time"

type Delivery struct {
	ID             uint64         `gorm:"primaryKey"`
	OrderID        uint64         `gorm:"not null"`
	Provider       string         `gorm:"not null"`
	Origin         Location       `gorm:"embedded;embeddedPrefix:origin_"`
	Destination    Location       `gorm:"embedded;embeddedPrefix:destination_"`
	TimeFrameStart time.Time      `gorm:"not null"`
	TimeFrameEnd   time.Time      `gorm:"not null"`
	Status         DeliveryStatus `gorm:"type:delivery_status;not null"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
}
