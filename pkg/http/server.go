package http

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
			// Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}




// func TestDispatcher(t *testing.T) {
// 	logger := zerolog.Nop()
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	t.Run("successful order processing", func(t *testing.T) {
// 		repo := new(MockRepository)
// 		tplClient := new(MockTPLClient)
		
// 		orders := []models.Order{
// 			{
// 				ID:       "order-1",
// 				Deadline: time.Now().Add(1 * time.Hour),
// 			},
// 		}

// 		// Mock expectations
// 		repo.On("GetPendingOrders", mock.Anything).Return(orders, nil).Once()
// 		tplClient.On("CreateShipment", mock.Anything, "order-1").Return(nil).Once()
// 		repo.On("MarkOrderCompleted", mock.Anything, "order-1").Return(nil).Once()

// 		dispatcher := scheduler.NewDispatcher(
// 			repo,
// 			tplClient,
// 			logger,
// 			1*time.Hour, // check interval
// 			3,           // max retries
// 		)

// 		// Run processing once
// 		dispatcher.ProcessOrders(ctx)

// 		repo.AssertExpectations(t)
// 		tplClient.AssertExpectations(t)
// 	})

// 	t.Run("retry failed shipments", func(t *testing.T) {
// 		repo := new(MockRepository)
// 		tplClient := new(MockTPLClient)
		
// 		orders := []models.Order{
// 			{
// 				ID:       "order-1",
// 				Deadline: time.Now().Add(1 * time.Hour),
// 			},
// 		}

// 		// Mock expectations
// 		repo.On("GetPendingOrders", mock.Anything).Return(orders, nil).Once()
// 		tplClient.On("CreateShipment", mock.Anything, "order-1").
// 			Return(errors.New("connection failed")).
// 			Times(3) // Expect 3 retries
// 		// No MarkOrderCompleted call expected

// 		dispatcher := scheduler.NewDispatcher(
// 			repo,
// 			tplClient,
// 			logger,
// 			1*time.Hour,
// 			3,
// 		)

// 		dispatcher.ProcessOrders(ctx)

// 		repo.AssertExpectations(t)
// 		tplClient.AssertExpectations(t)
// 	})

// 	t.Run("expired orders are skipped", func(t *testing.T) {
// 		repo := new(MockRepository)
// 		tplClient := new(MockTPLClient)
		
// 		orders := []models.Order{
// 			{
// 				ID:       "order-1",
// 				Deadline: time.Now().Add(-1 * time.Hour),
// 			},
// 		}

// 		repo.On("GetPendingOrders", mock.Anything).Return(orders, nil).Once()

// 		dispatcher := scheduler.NewDispatcher(
// 			repo,
// 			tplClient,
// 			logger,
// 			1*time.Hour,
// 			3,
// 		)

// 		dispatcher.ProcessOrders(ctx)

// 		repo.AssertExpectations(t)
// 		tplClient.AssertNotCalled(t, "CreateShipment")
// 	})

// 	t.Run("context cancellation during processing", func(t *testing.T) {
// 		repo := new(MockRepository)
// 		tplClient := new(MockTPLClient)
		
// 		orders := []models.Order{
// 			{ID: "order-1", Deadline: time.Now().Add(1 * time.Hour)},
// 			{ID: "order-2", Deadline: time.Now().Add(1 * time.Hour)},
// 		}

// 		// Cancel context after first order
// 		ctx, cancel := context.WithCancel(context.Background())
// 		defer cancel()

// 		repo.On("GetPendingOrders", ctx).Return(orders, nil).Once()
// 		tplClient.On("CreateShipment", ctx, "order-1").Return(nil).Run(func(args mock.Arguments) {
// 			cancel()
// 		})
// 		repo.On("MarkOrderCompleted", ctx, "order-1").Return(nil).Once()

// 		dispatcher := scheduler.NewDispatcher(
// 			repo,
// 			tplClient,
// 			logger,
// 			1*time.Hour,
// 			3,
// 		)

// 		dispatcher.ProcessOrders(ctx)

// 		repo.AssertExpectations(t)
// 		tplClient.AssertExpectations(t)
// 	})

// 	t.Run("database error handling", func(t *testing.T) {
// 		repo := new(MockRepository)
// 		tplClient := new(MockTPLClient)
		
// 		repo.On("GetPendingOrders", mock.Anything).
// 			Return([]models.Order{}, errors.New("database connection failed")).
// 			Once()

// 		dispatcher := scheduler.NewDispatcher(
// 			repo,
// 			tplClient,
// 			logger,
// 			1*time.Hour,
// 			3,
// 		)

// 		dispatcher.ProcessOrders(ctx)

// 		repo.AssertExpectations(t)
// 	})
// }
