package dto

type CreateBillingRequest struct {
	CustomerID   uint `json:"customerId"`
	LoanID       uint `json:"loanId"`
	LoanAmount   int  `json:"loanAmount"`
	LoanInterest int  `json:"loanInterest"`
	LoanWeeks    int  `json:"loanWeeks"`
}

type GetBillingRequest struct {
	CustomerID uint `json:"customerId"`
	LoanID     uint `json:"loanId"`
}

type GetOutstandingRequest struct {
	CustomerID uint `json:"customerId"`
	LoanID     uint `json:"loanId"`
}
