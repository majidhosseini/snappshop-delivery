package entity

import "time"

type DeliveryAudit struct {
	ID         uint64         `gorm:"primaryKey"`
	DeliveryId uint64         `gorm:"not null"`
	Provider   string         `gorm:"not null"`
	Status     DeliveryStatus `gorm:"type:delivery_status;not null"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
}
