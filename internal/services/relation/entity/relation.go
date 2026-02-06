package entity

import (
	"appsite-go/internal/core/model"
)

// Relation Relation Entity
type Relation struct {
	model.Base
	ItemID       string `json:"item_id" gorm:"type:varchar(36);not null;index"`
	ItemType     string `json:"item_type" gorm:"type:varchar(32);not null;index"`
	RelationID   string `json:"relation_id" gorm:"type:varchar(36);not null;index"`
	RelationType string `json:"relation_type" gorm:"type:varchar(32);not null;index"`
	Type         string `json:"type" gorm:"type:varchar(16);index"` // e.g., follow, like, favorite
	Rate         int    `json:"rate" gorm:"default:0"`
	Status       string `json:"status" gorm:"type:varchar(32);default:'enabled'"`
}

// TableName table name
func (Relation) TableName() string {
	return "relation_combine"
}
