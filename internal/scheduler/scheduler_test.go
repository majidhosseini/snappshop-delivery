package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"snappshop.ir/internal/domain/entity"
	"snappshop.ir/internal/scheduler"
	"snappshop.ir/internal/tpl"
	"snappshop.ir/pkg/logger"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByTimeToDeliver(ctx context.Context) ([]entity.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockRepository) MarkOrderAsDispatched(ctx context.Context, orderNumber string) error {
	args := m.Called(ctx, orderNumber)
	return args.Error(0)
}

func (m *MockRepository) MarkOrderAsCompleted(ctx context.Context, orderNumber string) error {
	args := m.Called(ctx, orderNumber)
	return args.Error(0)
}

type MockTPLClient struct {
	mock.Mock
}

func (m *MockTPLClient) CreateShipment(ctx context.Context, order entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestDispatcher(t *testing.T) {
	logger := zerolog.Nop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("successful order processing", func(t *testing.T) {
		repo := new(MockRepository)
		tplClient := new(MockTPLClient)

		orders := []entity.Order{
			{
				ID:       123,
				OrderNumber: "order-1",
				Status:  entity.StatusCreated,
				TimeFrame: entity.TimeFrame{
					Start: time.Now(),
					End:   time.Now().Add(1 * time.Hour),
				},
				UserInfo: entity.UserInfo{
					Name:     "John Doe",
					Phone:    "1234567890",
					Address:  "123 Main St",
					Username: "johndoe",
					Email:    "johndoe@email.com",
				},
				Origin: entity.Location{
					Latitude:  37.7749,
					Longitude: 122.4194,
				},
				Destination: entity.Location{
					Latitude:  37.7749,
					Longitude: 122.4194,
				},
			},
			}

		// Mock expectations
		repo.On("GetByTimeToDeliver", mock.Anything).Return(orders, nil).Once()
		tplClient.On("CreateShipment", mock.Anything, orders[0]).Return(nil).Once()
		repo.On("MarkOrderAsCompleted", mock.Anything, "order-1").Return(nil).Once()

		dispatcher := scheduler.NewDispatcher(
			repo,
			tplClient,
			logger,
			1*time.Hour, // check interval
			3,
		)

		// Run processing once
		dispatcher.(ctx)

		repo.AssertExpectations(t)
		tplClient.AssertExpectations(t)

	})

	repo := new(MockRepository)
	tplClient := new(MockTPLClient)
	logger := logger.New("testing_service")
	checkInterval := 5 * time.Second
	maxRetries := 3

	dispatcher := scheduler.NewDispatcher(repo, tplClient, logger, checkInterval, maxRetries)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo.On("GetPendingOrders", ctx).Return([]entity.Order{}, nil)

	dispatcher.Start(ctx)

	repo.AssertExpectations(t)
	tplClient.AssertExpectations(t)
}
