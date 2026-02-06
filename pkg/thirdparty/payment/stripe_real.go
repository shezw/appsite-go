// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/refund"
	"github.com/stripe/stripe-go/v72/webhook"
)

// StripeReal implements real Stripe integration
type StripeReal struct {
	SecretKey     string
	WebhookSecret string
}

func NewStripeReal(secretKey, webhookSecret string) *StripeReal {
	stripe.Key = secretKey
	return &StripeReal{
		SecretKey:     secretKey,
		WebhookSecret: webhookSecret,
	}
}

func (g *StripeReal) Pay(req *Request) (*Response, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
		Description: stripe.String(req.Description),
	}
	// Add metadata
	params.AddMetadata("order_id", req.OrderID)

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe create pi failed: %w", err)
	}

	return &Response{
		TransactionID: pi.ID,
		// No direct PayURL for PI unless using Checkout Sessions. 
		// We return Data for frontend SDK to handle 'confirm'.
		Data: map[string]string{
			"client_secret": pi.ClientSecret, 
			"public_key":    "pk_...", // In real app, PK is config in frontend
		},
	}, nil
}

func (g *StripeReal) Verify(r *http.Request) (string, string, error) {
	// Read body
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(nil, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return "", "", fmt.Errorf("read body failed: %w", err)
	}

	// Verify signature
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), g.WebhookSecret)
	if err != nil {
		return "", "", ErrSignature
	}

	// Handle event
	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &pi)
		if err != nil {
			return "", "", err
		}
		
		orderID := pi.Metadata["order_id"]
		return orderID, "paid", nil

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		json.Unmarshal(event.Data.Raw, &pi)
		orderID := pi.Metadata["order_id"]
		return orderID, "closed", nil // or failed
	}

	return "", "", nil // Ignore other events
}

func (g *StripeReal) Refund(orderID string, amount int64, reason string) error {
	// Here orderID should be the Transaction ID (pi_...) stored in Order struct
	// If caller passes our internal OrderID, we can't find it unless we query DB.
	// Assume orderID passed here IS the transaction ID or we stored it.
	
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(orderID),
		Amount:        stripe.Int64(amount),
	}
	_, err := refund.New(params)
	return err
}
