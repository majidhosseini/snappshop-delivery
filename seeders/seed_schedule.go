package seeders

import (
	"context"
	"fmt"
	"time"

	"snappshop.ir/internal/delivery"
)

func SeedScheduledRequests(service *delivery.Service, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		req := &delivery.Request{
			OrderID:           time.Now().Format("SCH-20060102150405"),
			UserInfo:          delivery.UserInfo{Name: "ScheduledUser", Phone: "9999999999"},
			FromLoc:           [2]float64{37.7749, -122.4194},
			ToLoc:             [2]float64{34.0522, -118.2437},
			DeliveryTimeFrame: [2]string{"15:00", "17:00"},
			State:             string(delivery.StateInit),
		}

		if err := service.CreateRequest(ctx, req); err != nil {
			fmt.Printf("Failed to seed scheduled request: %v\n", err)
		} else {
			fmt.Printf("Scheduled request created: %s\n", req.OrderID)
		}
	}
}
