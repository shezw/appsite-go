package entity

import (
	"appsite-go/internal/core/model"
	"appsite-go/pkg/dbs"
)

// Request Form Request Entity
type Request struct {
	model.Base
	SaasID     string  `json:"saas_id" gorm:"type:varchar(36);index"`
	UserID     string  `json:"user_id" gorm:"type:varchar(36);index"`
	ItemID     string  `json:"item_id" gorm:"type:varchar(36);index"`
	ItemType   string  `json:"item_type" gorm:"type:varchar(32);index"`
	Open       bool    `json:"open" gorm:"default:false"`
	Form       dbs.Map `json:"form" gorm:"type:json"`
	Status     string  `json:"status" gorm:"type:varchar(12);default:'pending';index"` // pending, applied, rejected
	Expire     int64   `json:"expire" gorm:"default:0"`
	ApplyCall  dbs.Map `json:"apply_call" gorm:"type:json"`
	RejectCall dbs.Map `json:"reject_call" gorm:"type:json"`
}

// TableName table name
func (Request) TableName() string {
	return "form_request"
}
