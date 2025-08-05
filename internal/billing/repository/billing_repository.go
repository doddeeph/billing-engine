package repository

import (
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"gorm.io/gorm"
)

type BillingRepository interface {
	WithTransaction(tx *gorm.DB) BillingRepository
	Create(billing *model.Billing) error
	FindByID(ID uint) (*model.Billing, error)
	FindByCustomerIdAndLoanId(customerID, loanID uint, includePayments bool) (*model.Billing, error)
	UpdateOutstandingBalance(billingID uint, balance int) error
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

func (r *billingRepository) Create(billing *model.Billing) error {
	return r.db.Create(billing).Error
}

func (r *billingRepository) FindByID(ID uint) (*model.Billing, error) {
	var billing model.Billing
	if err := r.db.Preload("Payments").First(&billing, ID).Error; err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *billingRepository) FindByCustomerIdAndLoanId(customerID, loanID uint, includePayments bool) (*model.Billing, error) {
	var billing model.Billing
	query := r.db.Model(&model.Billing{})
	if includePayments {
		query = query.Preload("Payments")
	}
	err := query.Where("customer_id = ? AND loan_id = ?", customerID, loanID).First(&billing).Error
	if err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *billingRepository) UpdateOutstandingBalance(billingID uint, balance int) error {
	return r.db.Model(&model.Billing{}).Where("id = ?", billingID).Update("outstanding_balance", balance).Error
}
