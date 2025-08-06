package model

type Billing struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CustomerID   uint      `gorm:"not null" json:"customerId"`
	LoanID       uint      `gorm:"uniqueIndex:idx_loan_id;not null" json:"loanId"`
	LoanAmount   int       `gorm:"not null" json:"loanAmount"`
	LoanWeeks    int       `gorm:"not null" json:"loanWeeks"`
	LoanInterest int       `gorm:"not null" json:"loanInterest"`
	Outstanding  int       `gorm:"not null" json:"outstanding"`
	Payments     []Payment `gorm:"foreignKey:BillingID"`
	CommonModel
}
