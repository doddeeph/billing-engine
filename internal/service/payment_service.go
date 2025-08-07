package service

import (
	"context"
	"fmt"
	"time"

	"github.com/doddeeph/billing-engine/internal/dto"
	"github.com/doddeeph/billing-engine/internal/repository"
	"gorm.io/gorm"
)

type PaymentService interface {
	MakePayment(ctx context.Context, billingId uint, req dto.PaymentRequest) (*dto.PaymentResponse, error)
}

type paymentServiceImpl struct {
	repo       repository.PaymentRepository
	billingSvc BillingService
}

func NewPaymentService(repo repository.PaymentRepository, billingSvc BillingService) PaymentService {
	return &paymentServiceImpl{repo: repo, billingSvc: billingSvc}
}

func (svc *paymentServiceImpl) MakePayment(ctx context.Context, billingId uint, req dto.PaymentRequest) (*dto.PaymentResponse, error) {
	var paymentResp *dto.PaymentResponse
	err := svc.repo.WithDB().Transaction(func(trx *gorm.DB) error {
		trxBillingSvc := svc.billingSvc.WithTransaction(trx)
		trxPaymentRepo := svc.repo.WithTransaction(trx)

		billing, err := trxBillingSvc.GetBilling(ctx, billingId)
		if err != nil {
			return err
		}

		if req.Week < 0 || req.Week > billing.LoanWeeks {
			return fmt.Errorf("Payment is outside %d loan week.", billing.LoanWeeks)
		}

		payment, err := trxPaymentRepo.FindByBillingIdAndWeek(ctx, billing.ID, req.Week)
		if err != nil {
			return err
		}
		if payment.Paid {
			return fmt.Errorf("Week %d has been paid.", req.Week)
		}
		if req.Amount < payment.Amount {
			return fmt.Errorf("Insufficient loan amount paid for week %d", req.Week)
		}

		payment.Paid = true
		now := time.Now()
		payment.PaidDate = &now
		updatedPayment, err := trxPaymentRepo.UpdatePaid(ctx, payment)
		if err != nil {
			return err
		}

		updatedOutstanding := billing.Outstanding - payment.Amount
		err = trxBillingSvc.UpdateOutstanding(ctx, billing.ID, updatedOutstanding)
		if err != nil {
			return err
		}

		paymentResp = &dto.PaymentResponse{
			CustomerID:  billing.CustomerID,
			LoanID:      billing.LoanID,
			Outstanding: updatedOutstanding,
			Payment:     *updatedPayment,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paymentResp, nil
}
