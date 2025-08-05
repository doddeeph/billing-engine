package model

import "gorm.io/gorm"

type Payment struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	BillingID uint `json:"billingId"`
	Amount    int  `json:"amount"`
	Week      int  `json:"week"`
	Paid      bool `json:"paid"`
	gorm.Model
}
