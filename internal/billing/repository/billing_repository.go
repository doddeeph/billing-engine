package repository

import (
	"context"

	"github.com/doddeeph/billing-engine/internal/billing/model"
	"gorm.io/gorm"
)

type BillingRepository interface {
	WithTransaction(tx *gorm.DB) BillingRepository
	Create(ctx context.Context, billing *model.Billing) error
	FindByID(ctx context.Context, ID uint) (*model.Billing, error)
	FindByCustomerIdAndLoanId(ctx context.Context, customerID, loanID uint, includePayments bool) (*model.Billing, error)
	UpdateOutstanding(ctx context.Context, billingID uint, balance int) error
}

type billingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) BillingRepository {
	return &billingRepository{db}
}

func (r *billingRepository) WithTransaction(tx *gorm.DB) BillingRepository {
	return &billingRepository{tx}
}

func (r *billingRepository) Create(ctx context.Context, billing *model.Billing) error {
	return r.db.WithContext(ctx).Create(billing).Error
}

func (r *billingRepository) FindByID(ctx context.Context, ID uint) (*model.Billing, error) {
	var billing model.Billing
	if err := r.db.WithContext(ctx).Preload("Payments").First(&billing, ID).Error; err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *billingRepository) FindByCustomerIdAndLoanId(ctx context.Context, customerID, loanID uint, includePayments bool) (*model.Billing, error) {
	var billing model.Billing
	query := r.db.WithContext(ctx).Model(&model.Billing{})
	if includePayments {
		query = query.Preload("Payments")
	}
	err := query.Where("customer_id = ? AND loan_id = ?", customerID, loanID).First(&billing).Error
	if err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *billingRepository) UpdateOutstanding(ctx context.Context, billingID uint, balance int) error {
	return r.db.WithContext(ctx).Model(&model.Billing{}).Where("id = ?", billingID).Update("outstanding", balance).Error
}
