package delivery

import (
	"errors"
	"time"

	"snappshop.ir/internal/domain/entity"
)

type Request struct {
	OrderNumber    string    `json:"orderNumber"`
	UserInfo       UserInfo  `json:"userInfo"`
	Origin         Location  `json:"origin"`
	Destination    Location  `json:"destination"`
	StartTimeFrame time.Time `json:"startTimeFrame"`
	EndTimeFrame   time.Time `json:"endTimeFrame"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (r *Request) Validate() error {
	if r.OrderNumber == "" || (r.Origin == Location{}) || (r.Destination == Location{}) {
		return errors.New("invalid request")
	}
	return nil
}

func (r *Request) toOrder() *entity.Order {
	return &entity.Order{
		OrderNumber: r.OrderNumber,
		UserInfo:    *r.UserInfo.toEntity(),
		Status:      entity.StatusCreated,
		Origin:      entity.Location{Latitude: r.Origin.Lat, Longitude: r.Origin.Lng},
		Destination: entity.Location{Latitude: r.Destination.Lat, Longitude: r.Destination.Lng},
		TimeFrame:   entity.TimeFrame{Start: r.StartTimeFrame, End: r.EndTimeFrame},
	}
}

type UserInfo struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (ui *UserInfo) toEntity() *entity.UserInfo {
	return &entity.UserInfo{
		Name:     ui.Name,
		Phone:    ui.Phone,
		Address:  ui.Address,
		Username: ui.Username,
		Email:    ui.Email,
	}
}
