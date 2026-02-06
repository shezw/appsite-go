package payment_test

import (
	"net/http"
	"net/url"
	"testing"

	"appsite-go/pkg/thirdparty/payment"
)

func TestMockGateway(t *testing.T) {
	gw := &payment.MockGateway{}

	// Pay
	req := &payment.Request{OrderID: "ORD-123", Amount: 100}
	resp, err := gw.Pay(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.TransactionID != "mock-ORD-123" {
		t.Error("Wrong tx ID")
	}

	// Verify
	u, _ := url.Parse("http://cb.com?order_id=ORD-123")
	r := &http.Request{URL: u}
	oid, status, err := gw.Verify(r)
	if err != nil {
		t.Fatal(err)
	}
	if oid != "ORD-123" || status != "paid" {
		t.Error("Verify failed")
	}

	// Refund
	if err := gw.Refund("ORD-123", 100, "reason"); err != nil {
		t.Fatal(err)
	}
}
