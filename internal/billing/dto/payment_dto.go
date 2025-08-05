package dto

type MakePaymetRequest struct {
	CustomerID uint `json:"customerId"`
	LoanID     uint `json:"loanId"`
	Week       int  `json:"week"`
}
