package entity

import "time"

type DeliveryStatus string

const (
	StateInit      DeliveryStatus = "init"
	StateIsFinding DeliveryStatus = "isFinding"
	StateFound     DeliveryStatus = "found"
	StateNotFound  DeliveryStatus = "notFound"
	StateDelivered DeliveryStatus = "delivered"
)

type Location struct {
	Latitude  float64 `gorm:"not null"`
	Longitude float64 `gorm:"not null"`
}

type TimeFrame struct {
	Days  int       `gorm:"not null"` // Number of days later
	Start time.Time `gorm:"not null"` // Start time
	End   time.Time `gorm:"not null"` // End time
}

type UserInfo struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
