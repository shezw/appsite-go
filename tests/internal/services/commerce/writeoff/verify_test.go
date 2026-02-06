package writeoff_test

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/commerce/coupon"
	"appsite-go/internal/services/commerce/entity"
	"appsite-go/internal/services/commerce/writeoff"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	// entity.UserCoupon/Coupon are needed
	db.AutoMigrate(&entity.Coupon{}, &entity.UserCoupon{})
	return db
}

func TestWriteOff(t *testing.T) {
	db := setupDB(t)
	cSvc := coupon.NewService(db)
	wSvc := writeoff.NewService(cSvc)

	// Setup: Create a coupon and issue it
	c := &entity.Coupon{
		Title: "Test Coupon",
		TotalCount: 10,
		Status: "enabled",
		StartTime: time.Now().Unix() - 100,
		EndTime: time.Now().Unix() + 100,
	}
	cSvc.CreateCoupon(c)

	uc, err := cSvc.Issue(func() string { return "u1" }, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	// 1. Verify Good Code
	ucFetch, cFetch, err := wSvc.VerifyCouponCode(uc.ID)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if ucFetch.ID != uc.ID {
		t.Error("ID mismatch")
	}
	if cFetch.ID != c.ID {
		t.Error("Coupon ID mismatch")
	}

	// 2. Write Off
	if err := wSvc.WriteOffCoupon(uc.ID, "staff-007"); err != nil {
		t.Fatalf("WriteOff failed: %v", err)
	}

	// 3. Verify Bad Code (Used)
	_, _, err = wSvc.VerifyCouponCode(uc.ID)
	if err != writeoff.ErrInvalidCode {
		t.Errorf("Expected ErrInvalidCode for used coupon, got %v", err)
	}
}
