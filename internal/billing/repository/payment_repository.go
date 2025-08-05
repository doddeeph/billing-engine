package repository

import (
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	WithTransaction(trx *gorm.DB) PaymentRepository
	WithDB() *gorm.DB
	FindByBillingIdAndWeek(billingID uint, week int) (*model.Payment, error)
	UpdatePaid(paymentID uint) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) WithTransaction(trx *gorm.DB) PaymentRepository {
	return &paymentRepository{trx}
}

func (r *paymentRepository) WithDB() *gorm.DB {
	return r.db
}

func (r *paymentRepository) FindByBillingIdAndWeek(billingID uint, week int) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("billing_id = ? AND week = ?", billingID, week).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePaid(paymentID uint) error {
	return r.db.Model(&model.Payment{}).Where("id = ?", paymentID).Update("paid", true).Error
}
