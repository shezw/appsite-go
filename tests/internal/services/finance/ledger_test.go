package finance_test

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/finance"
	"appsite-go/internal/services/finance/entity"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	// Important: Fix SQLite concurrency for transactions
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)

	return db
}

func TestLedger(t *testing.T) {
	db := setupDB(t)
	svc := finance.NewLedgerService(db)
	if err := svc.Migrate(); err != nil {
		t.Fatal(err)
	}

	uid := "u-100"
	asset := "point"

	// 1. Initial State
	bal, err := svc.GetBalance(uid, asset)
	if err != nil {
		t.Fatal(err)
	}
	if bal != 0 {
		t.Errorf("Expected 0 balance, got %d", bal)
	}

	// 2. Deposit (Reward)
	if err := svc.RecordTransaction(uid, asset, 100, "reward", "order-1", "Test Reward"); err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}
	
	bal, _ = svc.GetBalance(uid, asset)
	if bal != 100 {
		t.Errorf("Expected 100, got %d", bal)
	}

	time.Sleep(1 * time.Second) // Ensure timestamp diff for ordering

	// 3. Withdraw (Pay)
	if err := svc.RecordTransaction(uid, asset, -30, "payment", "order-2", "Test Pay"); err != nil {
		t.Fatalf("Withdraw failed: %v", err)
	}
	
	bal, _ = svc.GetBalance(uid, asset)
	if bal != 70 {
		t.Errorf("Expected 70, got %d", bal)
	}

	// 4. Overdraft Check
	if err := svc.RecordTransaction(uid, asset, -80, "payment", "order-3", "Fail Pay"); err != finance.ErrInsufficientFunds {
		t.Errorf("Expected ErrInsufficientFunds, got %v", err)
	}

	// 5. Check History
	var deals []entity.Deal
	var total int64
	deals, total, err = svc.ListDeals(uid, asset, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 { // +100, -30. Failed one not recorded.
		t.Errorf("Expected 2 deals, got %d", total)
	}
	if len(deals) != 2 {
		t.Errorf("Expected 2 items, got %d", len(deals))
	}
	if deals[0].Amount != -30 { // Desc order
		t.Errorf("Expected latest डील -30, got %d", deals[0].Amount)
	}
}
