package model

import "gorm.io/gorm"

type Billing struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	CustomerID         uint      `json:"customerId"`
	LoanID             uint      `gorm:"uniqueIndex" json:"loanId"`
	LoanAmount         int       `json:"loanAmount"`
	LoanWeeks          int       `json:"loanWeeks"`
	LoanInterest       int       `json:"loanInterest"`
	LoanWeeklyAmount   int       `json:"loanWeeklyAmount"`
	OutstandingBalance int       `json:"outstandingBalance"`
	Payments           []Payment `gorm:"foreignKey:BillingID"`
	gorm.Model
}
