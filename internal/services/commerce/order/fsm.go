// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"errors"
	"fmt"
	"time"

	"appsite-go/internal/services/commerce/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Status constants
const (
	StatusPending  = "pending"
	StatusPaid     = "paid"
	StatusShipping = "shipping"
	StatusDone     = "done"
	StatusClosed   = "closed"
	StatusRefunded = "refunded"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidState  = errors.New("invalid order state transition")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Create places a new order. 
// Note: Stock deduction should be done in the same transaction by the caller or we can inject StockService here.
// For flexibility, this method just creates the order record. Caller orchestrates the transaction.
func (s *Service) Create(tx *gorm.DB, order *entity.Order, items []entity.OrderItem) error {
	if tx == nil {
		tx = s.db
	}

	return tx.Transaction(func(t *gorm.DB) error {
		// 1. Save Order
		if order.OrderNo == "" {
			// Generate simple OrderNo if missing (time-random)
			order.OrderNo = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		if err := t.Create(order).Error; err != nil {
			return err
		}

		// 2. Save Items
		for i := range items {
			items[i].OrderID = order.ID
			if err := t.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// transition updates status with FSM check
func (s *Service) transition(orderID string, targetStatus string, check func(current string) bool, updates map[string]interface{}) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var o entity.Order
		// Lock row
		if err := tx.Clauses(clauseLocking).First(&o, "id = ?", orderID).Error; err != nil {
			return err
		}

		if !check(o.Status) {
			return fmt.Errorf("%w: cannot go from %s to %s", ErrInvalidState, o.Status, targetStatus)
		}

		if updates == nil {
			updates = make(map[string]interface{})
		}
		updates["status"] = targetStatus

		return tx.Model(&o).Updates(updates).Error
	})
}

// Pay marks order as paid
func (s *Service) Pay(orderID string, transactionID string) error {
	return s.transition(orderID, StatusPaid, func(curr string) bool {
		return curr == StatusPending
	}, map[string]interface{}{
		"pay_method":     "simulated", // dynamic in real world
		"transaction_id": transactionID,
		"pay_amount":     0, // In real world update with actual paid
	})
}

// Ship marks order as shipping
func (s *Service) Ship(orderID string) error {
	return s.transition(orderID, StatusShipping, func(curr string) bool {
		return curr == StatusPaid
	}, nil)
}

// Confirm marks order as done (received)
func (s *Service) Confirm(orderID string) error {
	return s.transition(orderID, StatusDone, func(curr string) bool {
		return curr == StatusShipping
	}, nil)
}

// Cancel marks order as closed (only if pending)
func (s *Service) Cancel(orderID string) error {
	return s.transition(orderID, StatusClosed, func(curr string) bool {
		return curr == StatusPending
	}, nil)
}

var clauseLocking = clause.Locking{Strength: "UPDATE"}
