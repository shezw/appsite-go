package stock_test

import (
	"sync"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/commerce/entity"
	"appsite-go/internal/services/commerce/stock"
)

func setupDB(t *testing.T) *gorm.DB {
	// Enable busy timeout to handle concurrency in SQLite
	// Also limit open connections to avoid "database table is locked" in excessive concurrency test
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_busy_timeout=10000"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1) // Force serialization for SQLite tests to avoid lock errors

	// Need Table for SKU and Log
	db.AutoMigrate(&entity.SKU{}, &entity.StockLog{})
	return db
}

func TestInventory(t *testing.T) {
	db := setupDB(t)
	svc := stock.NewInventoryService(db)

	// 1. Setup SKU
	sku := &entity.SKU{
		Code:  "TEST-SKU",
		Price: 100,
		Stock: 100,
	}
	db.Create(sku)
	skuID := sku.ID

	// 2. Deduct Normal
	err := svc.Deduct(skuID, 10, "order_1")
	if err != nil {
		t.Fatalf("Failed to deduct: %v", err)
	}
	
	current, _ := svc.GetStock(skuID)
	if current != 90 {
		t.Errorf("Expected 90, got %d", current)
	}

	// 3. Deduct Excessive (Should Fail)
	err = svc.Deduct(skuID, 100, "order_2")
	if err != stock.ErrInsufficientStock {
		t.Errorf("Expected ErrInsufficientStock, got %v", err)
	}

	// 4. Restore
	err = svc.Restore(skuID, 5, "order_1_cancel_part")
	if err != nil {
		t.Fatalf("Failed to restore: %v", err)
	}
	
	current, _ = svc.GetStock(skuID)
	if current != 95 {
		t.Errorf("Expected 95, got %d", current)
	}
	
	// 5. Check Logs
	var logs []entity.StockLog
	if err := db.Where("sku_id = ?", skuID).Find(&logs).Error; err != nil {
		t.Fatalf("Failed to query logs: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("Expected 2 logs, got %d", len(logs))
	}
	// Log 1: Deduct 10
	if logs[0].Quantity != -10 || logs[0].OrderID != "order_1" {
		t.Error("Log 1 mismatch")
	}
	// Log 2: Restore 5
	if logs[1].Quantity != 5 || logs[1].Type != "cancel" {
		t.Error("Log 2 mismatch")
	}
}

func TestInventory_Concurrency(t *testing.T) {
	db := setupDB(t)
	svc := stock.NewInventoryService(db)

	sku := &entity.SKU{Stock: 100}
	db.Create(sku)
	
	// Simulate 20 concurrent threads trying to buy 6 items each.
	// Total demand = 120. Stock = 100.
	// Should succeed 16 times (96 items), fail 4 times.
	// Or succeed X times until stock < 6.
	// 100/6 = 16.66 -> 16 successes. Remaining stock 4.
	
	var wg sync.WaitGroup
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := svc.Deduct(sku.ID, 6, "concurrent_order")
			mu.Lock()
			if err == nil {
				successCount++
			} else {
				failCount++
			}
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	
	current, _ := svc.GetStock(sku.ID)
	
	if successCount != 16 {
		t.Errorf("Expected 16 successful txs, got %d", successCount)
	}
	if current != 4 {
		t.Errorf("Expected remaining stock 4, got %d", current)
	}
}
