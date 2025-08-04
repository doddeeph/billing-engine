package repository

import (
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"gorm.io/gorm"
)

type BillingRepository interface {
	Create(billing *model.Billing) error
}

type billingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) BillingRepository {
	return &billingRepository{db}
}

func (r *billingRepository) Create(billing *model.Billing) error {
	return r.db.Create(billing).Error
}
