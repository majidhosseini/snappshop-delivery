package entity

import "time"

type Order struct {
	ID          uint64      `gorm:"primaryKey"`
	OrderNumber string      `gorm:"unique;not null"`
	Status      OrderStatus `gorm:"type:order_status;not null"`
	UserInfo    UserInfo    `gorm:"type:json;not null"`
	Origin      Location    `gorm:"embedded;embeddedPrefix:origin_"`
	Destination Location    `gorm:"embedded;embeddedPrefix:destination_"`
	TimeFrame   TimeFrame   `gorm:"embedded;embeddedPrefix:time_frame_"`
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`
}

type OrderStatus string

const (
	StatusCreated    OrderStatus = "created"
	StatusInProgress OrderStatus = "in_progress"
	StatusCompleted  OrderStatus = "completed"
	StatusCanceled   OrderStatus = "canceled"
)
