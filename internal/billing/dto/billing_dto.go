package dto

type CreateBillingRequest struct {
	CustomerID   uint `json:"customerId"`
	LoanID       uint `json:"loanId"`
	LoanAmount   int  `json:"loanAmount"`
	LoanInterest int  `json:"loanInterest"`
	LoanWeeks    int  `json:"loanWeeks"`
}
