package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"snappshop.ir/internal/domain/entity"
	"snappshop.ir/internal/scheduler"
	"snappshop.ir/pkg/logger"
)

type MockRepository struct {
	mock.Mock
}

// Create implements repository.OrderRepository.
func (m *MockRepository) Create(order *entity.Order) error {
	panic("unimplemented")
}

// Delete implements repository.OrderRepository.
func (m *MockRepository) Delete(id uint64) error {
	panic("unimplemented")
}

// GetByID implements repository.OrderRepository.
func (m *MockRepository) GetByID(id uint64) (*entity.Order, error) {
	panic("unimplemented")
}

// GetByOrderNumber implements repository.OrderRepository.
func (m *MockRepository) GetByOrderNumber(orderNumber string) (*entity.Order, error) {
	panic("unimplemented")
}

// MarkOrderCompleted implements repository.OrderRepository.
func (m *MockRepository) MarkOrderCompleted(ctx context.Context, orderNumber string) error {
	args := m.Called(ctx, orderNumber)
	return args.Error(0)
}

// Update implements repository.OrderRepository.
func (m *MockRepository) Update(order *entity.Order) error {
	panic("unimplemented")
}

func (m *MockRepository) GetByTimeToDeliver(ctx context.Context) ([]entity.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockRepository) MarkOrderAsDispatched(ctx context.Context, orderNumber string) error {
	args := m.Called(ctx, orderNumber)
	return args.Error(0)
}

type MockTPLClient struct {
	mock.Mock
}

func (m *MockTPLClient) CreateShipment(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func TestDispatcher(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("successful order processing", func(t *testing.T) {
		repo := new(MockRepository)
		tplClient := new(MockTPLClient)
		// logger := zerolog.Nop()

		orders := []entity.Order{
			{
				ID:          123,
				OrderNumber: "order-1",
				Status:      entity.StatusCreated,
				TimeFrame: entity.TimeFrame{
					Start: time.Now().Add(1 * time.Hour),
					End:   time.Now().Add(3 * time.Hour),
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
		repo.On("GetByTimeToDeliver", mock.Anything).Return(orders, nil).Maybe()
		tplClient.On("CreateShipment", mock.Anything, "order-1").Return(nil).Maybe()
		repo.On("MarkOrderCompleted", mock.Anything, "order-1").Return(nil).Maybe()

		logger2 := logger.New("test")

		dispatcher := scheduler.NewDispatcher(
			repo,
			tplClient,
			logger2.Logger,
			50*time.Second, // check interval
			3,
		)

		// Run processing once
		dispatcher.Start(ctx)

		repo.AssertExpectations(t)
		tplClient.AssertExpectations(t)
	})

	// dispatcher := scheduler.NewDispatcher(repo, tplClient, logger, checkInterval, maxRetries)

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// repo.On("GetPendingOrders", ctx).Return([]entity.Order{}, nil)

	// dispatcher.Start(ctx)

	// repo.AssertExpectations(t)
	// tplClient.AssertExpectations(t)
}
