// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sms

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// AliyunSender implements Sender for Alibaba Cloud SMS (Dysmsapi)
type AliyunSender struct {
	Client *dysmsapi.Client
	
	SignName string // The exact Sign Name configured in console
}

// NewAliyunSender creates a client
// regionID: e.g. "cn-hangzhou"
func NewAliyunSender(regionID, accessKeyID, accessKeySecret, signName string) (*AliyunSender, error) {
	client, err := dysmsapi.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	return &AliyunSender{
		Client:   client,
		SignName: signName,
	}, nil
}

func (s *AliyunSender) Send(phone string, message string) error {
	// Aliyun generally does not allow sending arbitrary text unless using specific templates or international SMS.
	// We return error to enforce Template usage, or you could use a "Default" template here.
	return fmt.Errorf("Aliyun SMS requires template usage. Call SendTemplate instead")
}

func (s *AliyunSender) SendTemplate(phone string, templateCode string, params map[string]string) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = s.SignName
	request.TemplateCode = templateCode

	if params != nil && len(params) > 0 {
		paramBytes, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("failed to marshal params: %w", err)
		}
		request.TemplateParam = string(paramBytes)
	}

	response, err := s.Client.SendSms(request)
	if err != nil {
		return err
	}

	if response.Code != "OK" {
		return fmt.Errorf("aliyun sms error: %s - %s", response.Code, response.Message)
	}

	return nil
}
