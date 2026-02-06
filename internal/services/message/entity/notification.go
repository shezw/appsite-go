package entity

import (
	"appsite-go/internal/core/model"
	"appsite-go/pkg/dbs"
)

// Notification Notification Entity
type Notification struct {
	model.Base
	SaasID     string  `json:"saas_id" gorm:"type:varchar(36);index"`
	SenderID   string  `json:"sender_id" gorm:"type:varchar(36);index;not null"`
	ReceiverID string  `json:"receiver_id" gorm:"type:varchar(36);index;not null"`
	ReplyID    string  `json:"reply_id" gorm:"type:varchar(36)"`
	Type       string  `json:"type" gorm:"type:varchar(32);default:'normal'"` // message, notify, suggest
	Status     string  `json:"status" gorm:"type:varchar(24);default:'sent'"` // sent, read
	Content    string  `json:"content" gorm:"type:varchar(512)"`
	Link       string  `json:"link" gorm:"type:varchar(255)"`
	LinkParams dbs.Map `json:"link_params" gorm:"type:json"`
	LinkType   string  `json:"link_type" gorm:"type:varchar(16)"`
}

// TableName table name
func (Notification) TableName() string {
	return "message_notification"
}
