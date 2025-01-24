package repository

import (
	"errors"

	"gorm.io/gorm"
	"snappshop.ir/internal/domain/entity"
)

// DeliveryRepository defines the interface for the delivery repository
type DeliveryRepository interface {
	Create(delivery *entity.Delivery) error
	GetByID(id uint64) (*entity.Delivery, error)
	Update(delivery *entity.Delivery) error
	Delete(id uint64) error
}

// deliveryRepository is the implementation of DeliveryRepository
type deliveryRepository struct {
	db *gorm.DB
}

// NewDeliveryRepository creates a new instance of deliveryRepository
func NewDeliveryRepository(db *gorm.DB) DeliveryRepository {
	return &deliveryRepository{db: db}
}

// Create adds a new delivery to the database
func (r *deliveryRepository) Create(delivery *entity.Delivery) error {
	return r.db.Create(delivery).Error
}

// GetByID retrieves a delivery by its ID
func (r *deliveryRepository) GetByID(id uint64) (*entity.Delivery, error) {
	var delivery entity.Delivery
	if err := r.db.First(&delivery, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &delivery, nil
}

// Update modifies an existing delivery in the database
func (r *deliveryRepository) Update(delivery *entity.Delivery) error {
	return r.db.Save(delivery).Error
}

// Delete removes a delivery from the database by its ID
func (r *deliveryRepository) Delete(id uint64) error {
	return r.db.Delete(&entity.Delivery{}, id).Error
}
