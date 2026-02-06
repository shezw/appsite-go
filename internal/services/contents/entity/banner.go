package entity

import (
	"appsite-go/internal/core/model"
)

// Banner Banner Entity
type Banner struct {
	model.Base
	SaasID     string `json:"saas_id" gorm:"type:varchar(36);index"`
	Position   string `json:"position" gorm:"type:varchar(32);index"`
	Title      string `json:"title" gorm:"type:varchar(255);not null"`
	Cover      string `json:"cover" gorm:"type:varchar(255)"`
	Link       string `json:"link" gorm:"type:varchar(255)"`
	ClickTimes int    `json:"click_times" gorm:"default:0"`
	ViewTimes  int    `json:"view_times" gorm:"default:0"`
	Status     string `json:"status" gorm:"type:varchar(32);default:'enabled';index"`
	Featured   bool   `json:"featured" gorm:"default:false;index"`
	Sort       int    `json:"sort" gorm:"default:0;index"`
}

// TableName table name
func (Banner) TableName() string {
	return "item_banner"
}
