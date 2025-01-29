package seeders

// func SeedRequests(service *delivery.Service, count int) {
// 	ctx := context.Background()
// 	rand.Seed(time.Now().UnixNano())

// 	for i := 0; i < count; i++ {
// 		req := &delivery.Request{
// 			OrderID: fmt.Sprintf("ORD-%d", i+1),
// 			UserInfo: delivery.UserInfo{
// 				Name:  fmt.Sprintf("User%d", i+1),
// 				Phone: fmt.Sprintf("12345%d", i+1),
// 			},
// 			FromLoc:           [2]float64{rand.Float64() * 90, rand.Float64() * 180},
// 			ToLoc:             [2]float64{rand.Float64() * 90, rand.Float64() * 180},
// 			DeliveryTimeFrame: [2]string{"09:00", "11:00"},
// 			State:             string(delivery.StateInit),
// 		}

// 		if err := service.CreateRequest(ctx, req); err != nil {
// 			fmt.Printf("Failed to seed request %s: %v\n", req.OrderID, err)
// 		}
// 	}
// 	fmt.Printf("Seeded %d requests successfully.\n", count)
// }
