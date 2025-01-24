package delivery

import (
	"errors"
	"time"

	"snappshop.ir/internal/domain/entity"
)

type Request struct {
	OrderNumber    string     `json:"orderNumber"`
	UserInfo       UserInfo   `json:"userInfo"`
	FromLoc        [2]float64 `json:"fromLoc"`
	ToLoc          [2]float64 `json:"toLoc"`
	StartTimeFrame time.Time  `json:"startTimeFrame"`
	EndTimeFrame   time.Time  `json:"endTimeFrame"`
}

func (r *Request) Validate() error {
	if r.OrderNumber == "" || r.FromLoc == [2]float64{} || r.ToLoc == [2]float64{} {
		return errors.New("invalid request")
	}
	return nil
}

func (r *Request) toOrder() *entity.Order {
	return &entity.Order{
		OrderNumber: r.OrderNumber,
		UserInfo:    *r.UserInfo.toEntity(),
		Status:      entity.StatusCreated,
		Origin:      entity.Location{Latitude: r.FromLoc[0], Longitude: r.FromLoc[1]},
		Destination: entity.Location{Latitude: r.ToLoc[0], Longitude: r.ToLoc[1]},
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
