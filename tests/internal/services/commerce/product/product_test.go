package product_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"appsite-go/internal/services/commerce/entity"
	"appsite-go/internal/services/commerce/product"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestProduct_CRUD(t *testing.T) {
	db := setupDB(t)
	svc := product.NewService(db)

	prod := &entity.Product{
		Title: "Test Phone",
		Price: 99900,
		Status: "on_sale",
	}

	// 1. Create SPU
	if err := svc.CreateProduct(prod); err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	if prod.ID == "" {
		t.Error("Product ID missing")
	}

	// 2. Create SKUs
	sku1 := &entity.SKU{
		ProductID: prod.ID,
		Title: "Black 128GB",
		Code: "P-B-128",
		Price: 99900,
		Stock: 10,
	}
	if err := svc.CreateSKU(sku1); err != nil {
		t.Fatalf("Failed to create SKU: %v", err)
	}

	// 3. Get
	fetched, err := svc.GetProduct(prod.ID)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}
	if fetched.Title != "Test Phone" {
		t.Errorf("Title mismatch")
	}

	// 4. Update
	updates := map[string]interface{}{"price": 89900}
	svc.UpdateProduct(prod.ID, updates)
	svc.UpdateSKU(sku1.ID, updates)

	fetchedSKU, _ := svc.GetSKU(sku1.ID)
	if fetchedSKU.Price != 89900 {
		t.Error("SKU update failed")
	}

	// 5. List SKUs
	skus, err := svc.ListSkus(prod.ID)
	if err != nil {
		t.Fatalf("Failed to list SKUs: %v", err)
	}
	if len(skus) != 1 {
		t.Errorf("Expected 1 SKU, got %d", len(skus))
	}

	// 6. Delete
	if err := svc.DeleteProduct(prod.ID); err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}
	
	skus, _ = svc.ListSkus(prod.ID) // Ideally should be 0, but this depends on logical vs hard delete.
	// Model has Base but not SoftDelete by default unless I explicitly add it.
	// CRUD.Remove does hard delete if SoftDelete struct is not embedded or used?
	// Checking Base struct... it has DeletedAt?
	// Base struct in model/base.go:
	// type Base struct { ... }
	// It does NOT have DeletedAt.
	// SoftDelete is a separate struct.
	// Entity Product embeds Base.
	// So CRUD.Remove -> db.Delete -> hard delete.
	// So SKUs should be gone if transaction worked (we did manual delete).
	
	if len(skus) != 0 {
		t.Errorf("Expected 0 SKUs after delete, got %d", len(skus))
	}
}
