package dto

type CreateBillingRequest struct {
	CustomerID   uint `json:"customerId"`
	LoanID       uint `json:"loanId"`
	LoanAmount   int  `json:"loanAmount"`
	LoanInterest int  `json:"loanInterest"`
	LoanWeeks    int  `json:"loanWeeks"`
}

type OutstandingResponse struct {
	BillingID   uint `json:"billingId"`
	CustomerID  uint `json:"customerId"`
	LoanID      uint `json:"loanId"`
	Outstanding int  `json:"outstanding"`
}

type DelinquentResponse struct {
	BillingID    uint `json:"billingId"`
	CustomerID   uint `json:"customerId"`
	LoanID       uint `json:"loanId"`
	IsDelinquent bool `json:"isDelinquent"`
}
