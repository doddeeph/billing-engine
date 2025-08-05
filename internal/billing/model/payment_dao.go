package model

import "gorm.io/gorm"

type Payment struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	BillingID uint `json:"billingId"`
	Amount    int  `json:"amount"`
	Week      int  `json:"week"`
	Paid      bool `gorm:"default:false" json:"paid"`
	gorm.Model
}
