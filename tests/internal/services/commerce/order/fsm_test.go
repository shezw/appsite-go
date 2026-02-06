package order_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/commerce/entity"
	"appsite-go/internal/services/commerce/order"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&entity.Order{}, &entity.OrderItem{})
	return db
}

func TestOrderFSM(t *testing.T) {
	db := setupDB(t)
	svc := order.NewService(db)

	// 1. Create
	o := &entity.Order{
		UserID:      "u1",
		TotalAmount: 1000,
		PayAmount:   1000,
	}
	items := []entity.OrderItem{
		{
			ProductID: "p1",
			Title:     "Something",
			Price:     1000,
			Quantity:  1,
			Amount:    1000,
		},
	}
	
	if err := svc.Create(nil, o, items); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if o.Status != order.StatusPending {
		t.Errorf("Expected pending, got %s", o.Status)
	}

	// 2. Pay
	if err := svc.Pay(o.ID, "txn-123"); err != nil {
		t.Fatalf("Pay failed: %v", err)
	}
	
	// Check status
	var check entity.Order
	db.First(&check, "id = ?", o.ID)
	if check.Status != order.StatusPaid {
		t.Errorf("Expected paid, got %s", check.Status)
	}

	// 3. Double Pay (Should fail or be handled, here strict transition fails)
	if err := svc.Pay(o.ID, "txn-456"); err == nil {
		t.Error("Expected error on double pay, got nil")
	}

	// 4. Ship
	if err := svc.Ship(o.ID); err != nil {
		t.Fatalf("Ship failed: %v", err)
	}

	// 5. Cancel (Shipping -> Closed invalid)
	if err := svc.Cancel(o.ID); err == nil {
		t.Error("Expected error cancelling shipped order, got nil")
	}

	// 6. Confirm
	if err := svc.Confirm(o.ID); err != nil {
		t.Fatalf("Confirm failed: %v", err)
	}
	
	db.First(&check, "id = ?", o.ID)
	if check.Status != order.StatusDone {
		t.Errorf("Expected done, got %s", check.Status)
	}
}
