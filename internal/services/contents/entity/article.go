package entity

import (
	"appsite-go/internal/core/model"
	"appsite-go/pkg/dbs"
)

// Article Article Entity
type Article struct {
	model.Base
	CategoryID  string   `json:"category_id" gorm:"type:varchar(36);index"`
	AuthorID    string   `json:"author_id" gorm:"type:varchar(36);index"`
	AreaID      string   `json:"area_id" gorm:"type:varchar(36);index"`
	RegionID    string   `json:"region_id" gorm:"type:varchar(36);index"`
	SaasID      string   `json:"saas_id" gorm:"type:varchar(36);index"`
	Type        string   `json:"type" gorm:"type:varchar(32)"`
	Mode        string   `json:"mode" gorm:"type:varchar(32)"`
	Title       string   `json:"title" gorm:"type:varchar(255);not null;index"`
	Cover       string   `json:"cover" gorm:"type:varchar(255)"`
	Gallery     dbs.Slice       `json:"gallery" gorm:"type:json"`
	Attachments dbs.Slice       `json:"attachments" gorm:"type:json"`
	Video       string          `json:"video" gorm:"type:varchar(255)"`
	Link        string          `json:"link" gorm:"type:varchar(255)"`
	Tags        dbs.StringArray `json:"tags" gorm:"type:json"`
	Description string   `json:"description" gorm:"type:varchar(255)"`
	Introduce   string   `json:"introduce" gorm:"type:longtext"`
	ViewTimes   int      `json:"view_times" gorm:"default:0"`
	Status      string   `json:"status" gorm:"type:varchar(32);default:'enabled';index"`
	Featured    bool     `json:"featured" gorm:"default:false;index"`
	Sort        int      `json:"sort" gorm:"default:0;index"`
}

// TableName table name
func (Article) TableName() string {
	return "item_article"
}
