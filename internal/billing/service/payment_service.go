package service

import (
	"fmt"

	"github.com/doddeeph/billing-engine/internal/billing/dto"
	"github.com/doddeeph/billing-engine/internal/billing/repository"
	"gorm.io/gorm"
)

type PaymentService interface {
	MakePayment(req dto.PaymetRequest) error
}

type paymentServiceImpl struct {
	repo       repository.PaymentRepository
	billingSvc BillingService
}

func NewPaymentService(repo repository.PaymentRepository, billingSvc BillingService) PaymentService {
	return &paymentServiceImpl{repo: repo, billingSvc: billingSvc}
}

func (svc *paymentServiceImpl) MakePayment(req dto.PaymetRequest) error {
	return svc.repo.WithDB().Transaction(func(trx *gorm.DB) error {
		trxBillingSvc := svc.billingSvc.WithTransaction(trx)
		trxPaymentRepo := svc.repo.WithTransaction(trx)

		billing, err := trxBillingSvc.FindByCustomerIdAndLoanId(req.CustomerID, req.LoanID, true)
		if err != nil {
			return err
		}

		if req.Week < 0 || req.Week > billing.LoanWeeks {
			return fmt.Errorf("Payment is outside %d loan week.", billing.LoanWeeks)
		}

		payment, err := trxPaymentRepo.FindByBillingIdAndWeek(billing.ID, req.Week)
		if err != nil {
			return err
		}
		if payment.Paid {
			return fmt.Errorf("Week %d has been paid.", req.Week)
		}
		err = trxPaymentRepo.UpdatePaid(payment.ID)
		if err != nil {
			return err
		}

		newOutstandingBalance := billing.OutstandingBalance - payment.Amount
		err = trxBillingSvc.UpdateOutstandingBalance(billing.ID, newOutstandingBalance)
		if err != nil {
			return err
		}

		return nil
	})
}
