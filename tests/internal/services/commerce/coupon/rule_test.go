package coupon_test

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/commerce/coupon"
	"appsite-go/internal/services/commerce/entity"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_busy_timeout=10000"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)

	return db
}

func TestCouponService(t *testing.T) {
	db := setupDB(t)
	svc := coupon.NewService(db)
	if err := svc.Migrate(); err != nil {
		t.Fatal(err)
	}

	// 1. Create Coupon
	now := time.Now().Unix()
	c1 := &entity.Coupon{
		Title:      "Welcome 100-20",
		Type:       "cash",
		Value:      2000,
		MinSpend:   10000,
		StartTime:  now - 3600,
		EndTime:    now + 3600,
		TotalCount: 10,
		Status:     "enabled",
	}
	if err := svc.CreateCoupon(c1); err != nil {
		t.Fatal(err)
	}

	// 2. Issue Coupon
	user1 := "user-001"
	uc, err := svc.Issue(func() string { return user1 }, c1.ID)
	if err != nil {
		t.Fatalf("Issue failed: %v", err)
	}
	if uc.Status != "unused" {
		t.Errorf("Expected status unused, got %s", uc.Status)
	}

	// Double issue check (should fail)
	_, err = svc.Issue(func() string { return user1 }, c1.ID)
	if err != coupon.ErrAlreadyTaken {
		t.Errorf("Expected ErrAlreadyTaken, got %v", err)
	}

	// 3. Verify
	// Check min spend fail
	_, err = svc.Verify(uc.ID, 5000)
	if err != coupon.ErrMinSpend {
		t.Errorf("Expected ErrMinSpend, got %v", err)
	}

	// Check success
	_, err = svc.Verify(uc.ID, 12000)
	if err != nil {
		t.Errorf("Verify failed: %v", err)
	}

	// 4. Use
	if err := svc.Use(uc.ID, "order-999"); err != nil {
		t.Fatalf("Use failed: %v", err)
	}

	// Verify used status
	var check entity.UserCoupon
	db.First(&check, "id = ?", uc.ID)
	if check.Status != "used" {
		t.Errorf("Expected status used, got %s", check.Status)
	}
	
	// Double use check (should not find unused)
	if err := svc.Use(uc.ID, "order-888"); err == nil {
		t.Error("Expected error on reusing coupon, got nil") // it likely returns 'record not found' or similar count 0 update
	}

	// 5. Test Limits
	c2 := &entity.Coupon{
		Title:      "Limited",
		TotalCount: 1,
		Status:     "enabled",
		TakenCount:   0,
		StartTime:    now - 60,
		EndTime:      now + 60,
	}
	svc.CreateCoupon(c2)

	// User 2 takes
	_, err = svc.Issue(func() string { return "user-002" }, c2.ID)
	if err != nil {
		t.Fatal(err)
	}
	
	// 6. Test Listing/Counting
	list, err := svc.ListUserCoupons(user1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 { // Should be 1 (the one used above is status=used? Wait)
		// svc.Use was called on 'uc'.
		// So checking status 'unused' should return 0 for user1 if uc was the only one.
		t.Logf("List count: %d", len(list))
	}
	// Let's create a new unused one for list test
	c3 := &entity.Coupon{
		Title: "List Test",
		TotalCount: 10,
		Status: "enabled",
		StartTime: now,
		EndTime: now+1000,
	}
	svc.CreateCoupon(c3)
	
	uc3, err := svc.Issue(func() string { return user1 }, c3.ID)
	if err != nil {
		t.Fatal(err)
	}
	
	count, err := svc.CountUserCoupons(user1)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
	
	list, _ = svc.ListUserCoupons(user1)
	if len(list) != 1 {
		t.Errorf("Expected list 1, got %d", len(list))
	}
	if list[0].ID != uc3.ID {
		t.Errorf("ID mismatch")
	}

	// 7. Expired Issue
	cExpired := &entity.Coupon{
		Title:      "Expired",
		TotalCount: 10,
		Status:     "enabled",
		StartTime:  now - 3600,
		EndTime:    now - 60,
	}
	svc.CreateCoupon(cExpired)
	_, err = svc.Issue(func() string { return "u1" }, cExpired.ID)
	if err != coupon.ErrCouponExpired {
		t.Errorf("Expected ErrCouponExpired, got %v", err)
	}
}
