package message

import (
	"appsite-go/internal/core/model"
	"appsite-go/internal/services/message/entity"

	"gorm.io/gorm"
)

// NotificationService handles message operations
type NotificationService struct {
	db   *gorm.DB
	repo *model.CRUD[entity.Notification]
}

// NewNotificationService initializes the service
func NewNotificationService(db *gorm.DB) *NotificationService {
	if db != nil {
		_ = db.AutoMigrate(&entity.Notification{})
	}
	return &NotificationService{
		db:   db,
		repo: model.NewCRUD[entity.Notification](db),
	}
}

// Send sends a notification
func (s *NotificationService) Send(senderID, receiverID, content, msgType string) error {
	notification := &entity.Notification{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType,
		Status:     "sent",
	}
	res := s.repo.Add(notification)
	return res.Error
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(id string) error {
	res := s.repo.Update(id, map[string]interface{}{"status": "read"})
	return res.Error
}

// MarkAllAsRead marks all notifications for a receiver as read
func (s *NotificationService) MarkAllAsRead(receiverID string) error {
	return s.db.Model(&entity.Notification{}).
		Where("receiver_id = ? AND status = ?", receiverID, "sent").
		Update("status", "read").Error
}

// List returns notifications with filters
func (s *NotificationService) List(page, size int, filters map[string]interface{}) ([]entity.Notification, int64, error) {
	res := s.repo.List(&model.ListParams{
		Page:     page,
		PageSize: size,
		Filters:  filters,
		Sort:     "created_at desc",
	})

	if !res.Success {
		return nil, 0, res.Error
	}

	data := res.Data.(map[string]interface{})
	list := data["list"].([]entity.Notification)
	total := data["total"].(int64)

	return list, total, nil
}

// UnreadCount returns the count of unread notifications
func (s *NotificationService) UnreadCount(receiverID string) (int64, error) {
	var count int64
	err := s.db.Model(&entity.Notification{}).
		Where("receiver_id = ? AND status = ?", receiverID, "sent").
		Count(&count).Error
	return count, err
}
