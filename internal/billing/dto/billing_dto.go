package dto

type CreateBillingDTO struct {
	CustomerID   uint `json:"customerId"`
	LoanID       uint `json:"loanId"`
	LoanAmount   int  `json:"loanAmount"`
	LoanInterest int  `json:"loanInterest"`
	LoanWeeks    int  `json:"loanWeeks"`
}

type CreateBillingRequest struct {
	CreateBillingDTO
}

type CreateBillingResponse struct {
	BillingID   uint `json:"billingId"`
	Outstanding int  `json:"outstanding"`
	CreateBillingDTO
}

type BaseResponse struct {
	BillingID  uint `json:"billingId"`
	CustomerID uint `json:"customerId"`
	LoanID     uint `json:"loanId"`
}

type OutstandingResponse struct {
	BaseResponse
	Outstanding int `json:"outstanding"`
}

type DelinquentResponse struct {
	BaseResponse
	IsDelinquent bool `json:"isDelinquent"`
}
