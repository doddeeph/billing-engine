package dto

type PaymetRequest struct {
	CustomerID uint `json:"customerId"`
	LoanID     uint `json:"loanId"`
	Week       int  `json:"week"`
}

type IsDelinquentRequest struct {
	CustomerID uint `json:"customerId"`
	LoanID     uint `json:"loanId"`
}