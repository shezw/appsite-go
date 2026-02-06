package entity

import (
	"appsite-go/internal/core/model"
)

// Tag Tag Entity
type Tag struct {
	model.Base
	SaasID      string `json:"saas_id" gorm:"type:varchar(36);index"`
	AuthorID    string `json:"author_id" gorm:"type:varchar(36);index"`
	Type        string `json:"type" gorm:"type:varchar(32);index"`
	Title       string `json:"title" gorm:"type:varchar(64);not null"`
	Cover       string `json:"cover" gorm:"type:varchar(255)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	Status      string `json:"status" gorm:"type:varchar(32);default:'enabled';index"`
	Featured    bool   `json:"featured" gorm:"default:false;index"`
	Sort        int    `json:"sort" gorm:"default:0;index"`
}

// TableName table name
func (Tag) TableName() string {
	return "item_tag"
}
