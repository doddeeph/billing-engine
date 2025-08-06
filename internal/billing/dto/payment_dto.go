package dto

import "github.com/doddeeph/billing-engine/internal/billing/model"

type PaymentRequest struct {
	Week   int `json:"week"`
	Amount int `json:"amount"`
}

type PaymentResponse struct {
	CustomerID  uint          `json:"customerId"`
	LoanID      uint          `json:"loanId"`
	Outstanding int           `json:"outstanding"`
	Payment     model.Payment `json:"payment"`
}
