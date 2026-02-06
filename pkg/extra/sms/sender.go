// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sms

import (
	"encoding/json"
	"fmt"
	"log"
)

// Sender interface for sending SMS
type Sender interface {
	// Send sends raw text message
	Send(phone string, message string) error
	
	// SendTemplate sends a templated message (e.g. Aliyun SMS / Twilio Template)
	// params is a map or struct of variables
	SendTemplate(phone string, templateCode string, params map[string]string) error
}

// ConsoleSender logs SMS to stdout (for local dev)
type ConsoleSender struct {
	Prefix string
}

func (s *ConsoleSender) Send(phone string, message string) error {
	log.Printf("[%s] SMS to %s: %s", s.Prefix, phone, message)
	return nil
}

func (s *ConsoleSender) SendTemplate(phone string, templateCode string, params map[string]string) error {
	pBytes, _ := json.Marshal(params)
	log.Printf("[%s] Template SMS to %s | Tpl: %s | Params: %s", s.Prefix, phone, templateCode, string(pBytes))
	return nil
}

// MockSender for testing
type MockSender struct {
	LastPhone    string
	LastMessage  string
	LastTemplate string
	LastParams   map[string]string
}

func (s *MockSender) Send(phone string, message string) error {
	s.LastPhone = phone
	s.LastMessage = message
	return nil
}

func (s *MockSender) SendTemplate(phone string, templateCode string, params map[string]string) error {
	s.LastPhone = phone
	s.LastTemplate = templateCode
	s.LastParams = params
	return nil
}
