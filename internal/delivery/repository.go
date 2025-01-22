package delivery

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	SaveRequest(ctx context.Context, req *Request) error
	GetRequest(ctx context.Context, reqID string) (*Request, error)
	UpdateRequest(ctx context.Context, req *Request) error
}

type GORMRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) *GORMRepository {
	return &GORMRepository{db: db}
}

func (r *GORMRepository) SaveRequest(ctx context.Context, req *Request) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *GORMRepository) GetRequest(ctx context.Context, reqID string) (*Request, error) {
	var req Request
	err := r.db.WithContext(ctx).First(&req, "order_id = ?", reqID).Error
	return &req, err
}

func (r *GORMRepository) UpdateRequest(ctx context.Context, req *Request) error {
	return r.db.WithContext(ctx).Save(req).Error
}
