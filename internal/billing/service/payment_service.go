package service

import (
	"fmt"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/model"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
	"gorm.io/gorm"
)

type PaymentService interface {
	MakePayment(req dto.PaymetRequest) error
	MissedPayment(req dto.PaymetRequest) error
}

type paymentService struct {
	repo       repository.PaymentRepository
	billingSvc BillingService
}

func NewPaymentService(repo repository.PaymentRepository, billingSvc BillingService) PaymentService {
	return &paymentService{repo: repo, billingSvc: billingSvc}
}

func (svc *paymentService) MakePayment(req dto.PaymetRequest) error {
	return svc.repo.WithDB().Transaction(func(tx *gorm.DB) error {
		txBillingSvc := svc.billingSvc.WithTransaction(tx)
		txPaymentRepo := svc.repo.WithTransaction(tx)

		billing, err := txBillingSvc.FindByCustomerIdAndLoanId(req.CustomerID, req.LoanID, true)
		if err != nil {
			return err
		}

		if req.Week < 0 || req.Week > billing.LoanWeeks {
			return fmt.Errorf("Payment is outside %d loan week.", billing.LoanWeeks)
		}

		for _, p := range billing.Payments {
			if p.Week == req.Week && p.Paid {
				return fmt.Errorf("Week %d has been paid.", req.Week)
			}
		}

		payment := &model.Payment{
			BillingID: billing.ID,
			Amount:    billing.LoanWeeklyAmount,
			Week:      req.Week,
			Paid:      true,
		}

		err = txPaymentRepo.Create(payment)
		if err != nil {
			return err
		}

		newOutstandingBalance := billing.OutstandingBalance - payment.Amount
		err = txBillingSvc.UpdateOutstandingBalance(billing.ID, newOutstandingBalance)
		if err != nil {
			return err
		}

		return nil
	})
}

func (svc *paymentService) MissedPayment(req dto.PaymetRequest) error {
	return svc.repo.WithDB().Transaction(func(tx *gorm.DB) error {
		txBillingSvc := svc.billingSvc.WithTransaction(tx)
		txPaymentRepo := svc.repo.WithTransaction(tx)

		billing, err := txBillingSvc.FindByCustomerIdAndLoanId(req.CustomerID, req.LoanID, true)
		if err != nil {
			return err
		}

		if req.Week < 0 || req.Week > billing.LoanWeeks {
			return fmt.Errorf("Payment is outside %d loan week.", billing.LoanWeeks)
		}

		for _, p := range billing.Payments {
			if p.Week == req.Week && p.Paid {
				return fmt.Errorf("Week %d has been paid.", req.Week)
			}
		}

		payment := &model.Payment{
			BillingID: billing.ID,
			Amount:    0,
			Week:      req.Week,
			Paid:      false,
		}

		return  txPaymentRepo.Create(payment)
	})
}