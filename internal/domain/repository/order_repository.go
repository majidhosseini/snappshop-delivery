package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"snappshop.ir/internal/domain/entity"
)

// OrderRepository defines the interface for the order repository
// TODO 1: Add withContext
type OrderRepository interface {
	Create(order *entity.Order) error
	GetByID(id uint64) (*entity.Order, error)
	Update(order *entity.Order) error
	Delete(id uint64) error
	GetByTimeToDeliver() ([]entity.Order, error)
	GetByOrderNumber(orderNumber string) (*entity.Order, error)
}

// orderRepository is the implementation of OrderRepository
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of orderRepository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create adds a new order to the database
func (r *orderRepository) Create(order *entity.Order) error {
	return r.db.Create(order).Error
}

// GetByID retrieves an order by its ID
func (r *orderRepository) GetByID(id uint64) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetByTimeToDeliver retrieves orders that need to be delivered within the next hour
func (r *orderRepository) GetByTimeToDeliver() ([]entity.Order, error) {
	var orders []entity.Order
	oneHourLater := time.Now().Add(1 * time.Hour)

	if err := r.db.Where("status = ? AND time_frame_from <= ?", entity.StatusCreated, oneHourLater).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetByOrderNumber retrieves an order by its order number
func (r *orderRepository) GetByOrderNumber(orderNumber string) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.Where("order_number = ?", orderNumber).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// Update modifies an existing order in the database
func (r *orderRepository) Update(order *entity.Order) error {
	return r.db.Save(order).Error
}

// Delete removes an order from the database by its ID
func (r *orderRepository) Delete(id uint64) error {
	return r.db.Delete(&entity.Order{}, id).Error
}
