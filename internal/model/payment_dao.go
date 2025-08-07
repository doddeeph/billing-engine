package model

import "time"

type Payment struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	BillingID uint       `gorm:"index;not null" json:"billingId"`
	Amount    int        `gorm:"not null" json:"amount"`
	Week      int        `gorm:"not null" json:"week"`
	Paid      bool       `gorm:"default:false" json:"paid"`
	StartDate time.Time  `gorm:"not null" json:"startDate"`
	DueDate   time.Time  `gorm:"not null" json:"dueDate"`
	PaidDate  *time.Time `json:"paidDate"`
	CommonModel
}
