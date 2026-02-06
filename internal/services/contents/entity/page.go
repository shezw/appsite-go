package entity

import (
	"appsite-go/internal/core/model"
)

// Page Page Entity
type Page struct {
	model.Base
	Alias     string `json:"alias" gorm:"type:varchar(32);uniqueIndex"`
	SaasID    string `json:"saas_id" gorm:"type:varchar(36);index"`
	AuthorID  string `json:"author_id" gorm:"type:varchar(36);index"`
	Title     string `json:"title" gorm:"type:varchar(64);not null"`
	Cover     string `json:"cover" gorm:"type:varchar(255)"`
	Introduce string `json:"introduce" gorm:"type:longtext"`
	Status    string `json:"status" gorm:"type:varchar(32);default:'enabled';index"`
	ViewTimes int    `json:"view_times" gorm:"default:0"`
	Featured  bool   `json:"featured" gorm:"default:false;index"`
	Sort      int    `json:"sort" gorm:"default:0;index"`
}

// TableName table name
func (Page) TableName() string {
	return "item_page"
}
