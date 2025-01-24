package repository

import (
	"gorm.io/gorm"
	"snappshop.ir/internal/domain/entity"
)

// DeliveryAuditRepository defines the interface for the delivery audit repository
type DeliveryAuditRepository interface {
	Insert(audit *entity.DeliveryAudit) error
	GetByDeliveryId(deliveryId uint64) ([]entity.DeliveryAudit, error)
}

// deliveryAuditRepository is the implementation of DeliveryAuditRepository
type deliveryAuditRepository struct {
	db *gorm.DB
}

// NewDeliveryAuditRepository creates a new instance of deliveryAuditRepository
func NewDeliveryAuditRepository(db *gorm.DB) DeliveryAuditRepository {
	return &deliveryAuditRepository{db: db}
}

// Insert adds a new delivery audit to the database
func (r *deliveryAuditRepository) Insert(audit *entity.DeliveryAudit) error {
	return r.db.Create(audit).Error
}

// GetByDeliveryId retrieves delivery audits by DeliveryId
func (r *deliveryAuditRepository) GetByDeliveryId(deliveryId uint64) ([]entity.DeliveryAudit, error) {
	var audits []entity.DeliveryAudit
	if err := r.db.Where("delivery_id = ?", deliveryId).Find(&audits).Error; err != nil {
		return nil, err
	}
	return audits, nil
}
