package repository

import (
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	WithTransaction(tx *gorm.DB) PaymentRepository
	WithDB() *gorm.DB
	Create(payment *model.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) WithTransaction(tx *gorm.DB) PaymentRepository {
	return &paymentRepository{tx}
}

func (r *paymentRepository) WithDB() *gorm.DB {
	return r.db
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}
