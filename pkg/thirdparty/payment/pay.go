// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"errors"
	"net/http"
)

var (
	ErrPaymentFailed = errors.New("payment creation failed")
	ErrSignature     = errors.New("invalid signature")
)

type Request struct {
	OrderID     string
	Amount      int64 // In cents
	Currency    string
	Description string
	ClientIP    string
	NotifyURL   string
	ReturnURL   string
	Params      map[string]string // Extra params (e.g. openid for wechat)
}

type Response struct {
	TransactionID string            // Upstream ID (if available immediately)
	PayURL        string            // Redirect URL or QR code string
	Data          map[string]string // Form data or JSON for frontend SDK
}

// Gateway interface for payment providers
type Gateway interface {
	// Pay creates a payment order
	Pay(req *Request) (*Response, error)
	// Verify parses and verifies a callback/webhook request
	// returns: orderID, status(paid/closed), error
	Verify(r *http.Request) (string, string, error)
	// Refund refunds an order
	Refund(orderID string, amount int64, reason string) error
}

// MockGateway for testing
type MockGateway struct {
	MockID string
}

func (g *MockGateway) Pay(req *Request) (*Response, error) {
	return &Response{
		TransactionID: "mock-" + req.OrderID,
		PayURL:        "http://mock-payment.com/pay/" + req.OrderID,
	}, nil
}

func (g *MockGateway) Verify(r *http.Request) (string, string, error) {
	oid := r.URL.Query().Get("order_id")
	if oid == "" {
		return "", "", errors.New("missing order_id")
	}
	return oid, "paid", nil
}

func (g *MockGateway) Refund(orderID string, amount int64, reason string) error {
	return nil
}

// StripeGateway Skeleton
// Replaced by StripeReal in stripe_real.go
type StripeGateway struct {
	SecretKey string
}

