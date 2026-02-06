package entity

import (
	"appsite-go/internal/core/model"
	"appsite-go/pkg/dbs"
)

// Comment Comment Entity
type Comment struct {
	model.Base
	UserID   string   `json:"user_id" gorm:"type:varchar(36);index;not null"`
	ItemID   string   `json:"item_id" gorm:"type:varchar(36);index"`
	ItemType string   `json:"item_type" gorm:"type:varchar(32);index"`
	Title    string   `json:"title" gorm:"type:varchar(64)"`
	Content  string   `json:"content" gorm:"type:varchar(511)"`
	Details  dbs.Map  `json:"details" gorm:"type:json"`
	Status   string   `json:"status" gorm:"type:varchar(32);default:'enabled';index"`
	Featured bool     `json:"featured" gorm:"default:false;index"`
	Sort     int      `json:"sort" gorm:"default:0;index"`
}

// TableName table name
func (Comment) TableName() string {
	return "user_comment"
}
