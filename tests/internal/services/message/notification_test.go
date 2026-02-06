package message_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/message"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestNotification(t *testing.T) {
	db := setupDB(t)
	svc := message.NewNotificationService(db)

	sender := "user_1"
	receiver := "user_2"

	// 1. Send
	err := svc.Send(sender, receiver, "Hello World", "message")
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// 2. Unread Count
	count, err := svc.UnreadCount(receiver)
	if err != nil {
		t.Fatalf("Failed to count unread: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 unread, got %d", count)
	}

	// 3. List
	list, _, err := svc.List(1, 10, map[string]interface{}{"receiver_id": receiver})
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 message in list, got %d", len(list))
	}
	msg := list[0]
	if msg.Content != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", msg.Content)
	}

	// 4. Mark As Read
	err = svc.MarkAsRead(msg.ID)
	if err != nil {
		t.Fatalf("Failed to mark as read: %v", err)
	}

	count, _ = svc.UnreadCount(receiver)
	if count != 0 {
		t.Errorf("Expected 0 unread, got %d", count)
	}

	// 5. Mark All As Read
	svc.Send(sender, receiver, "Msg 2", "message")
	svc.Send(sender, receiver, "Msg 3", "message")
	
	err = svc.MarkAllAsRead(receiver)
	if err != nil {
		t.Fatalf("Failed to mark all as read: %v", err)
	}
	
	count, _ = svc.UnreadCount(receiver)
	if count != 0 {
		t.Errorf("Expected 0 unread, got %d", count)
	}
}
