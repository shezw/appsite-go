package entity

import (
	"appsite-go/internal/core/model"
)

// Word Shield Word Entity
type Word struct {
	model.Base
	Title    string `json:"title" gorm:"type:varchar(32);not null;index"`
	AuthorID string `json:"author_id" gorm:"type:varchar(36);index"`
	Type     string `json:"type" gorm:"type:varchar(16)"`
	Status   string `json:"status" gorm:"type:varchar(32);default:'enabled'"`
}

// TableName table name
func (Word) TableName() string {
	return "system_shieldword"
}
