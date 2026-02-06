package entity

import (
	"appsite-go/internal/core/model"
	"appsite-go/pkg/dbs"
)

// Media Media Entity
type Media struct {
	model.Base
	SaasID     string  `json:"saas_id" gorm:"type:varchar(36);index"`
	CategoryID string  `json:"category_id" gorm:"type:varchar(36);index"`
	AuthorID   string  `json:"author_id" gorm:"type:varchar(36);index"`
	Type       string  `json:"type" gorm:"type:varchar(16);index"` // image, video, file, etc.
	Server     int     `json:"server" gorm:"default:0"`            // 0: Local, 1: OSS
	URL        string  `json:"url" gorm:"type:varchar(255);not null"`
	Size       int64   `json:"size" gorm:"default:0"`
	Meta       dbs.Map `json:"meta" gorm:"type:json"`
	Password   string  `json:"password" gorm:"type:varchar(255)"`
	Status     string  `json:"status" gorm:"type:varchar(32);default:'enabled'"`
	Featured   bool    `json:"featured" gorm:"default:false;index"`
	Sort       int     `json:"sort" gorm:"default:0;index"`
}

// TableName table name
func (Media) TableName() string {
	return "item_media"
}
