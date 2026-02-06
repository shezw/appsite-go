package sms_test

import (
	"testing"

	"appsite-go/pkg/extra/sms"
)

func TestMockSender(t *testing.T) {
	s := &sms.MockSender{}
	
	// Test Send
	if err := s.Send("1234567890", "Hello"); err != nil {
		t.Fatal(err)
	}
	if s.LastPhone != "1234567890" {
		t.Errorf("Phone mismatch: %s", s.LastPhone)
	}
	if s.LastMessage != "Hello" {
		t.Errorf("Message mismatch: %s", s.LastMessage)
	}

	// Test Template
	params := map[string]string{"code": "1234"}
	if err := s.SendTemplate("0987654321", "SMS_1001", params); err != nil {
		t.Fatal(err)
	}
	if s.LastTemplate != "SMS_1001" {
		t.Errorf("Template mismatch")
	}
	if s.LastParams["code"] != "1234" {
		t.Errorf("Params mismatch")
	}
}

func TestConsoleSender(t *testing.T) {
	// Just coverage verify basically, acts as no-op test since it writes to log
	s := &sms.ConsoleSender{Prefix: "TEST"}
	s.Send("111", "msg")
	s.SendTemplate("111", "TPL", nil)
}
