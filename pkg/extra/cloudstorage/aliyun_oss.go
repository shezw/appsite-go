// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cloudstorage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSConfig holds configuration
type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	Domain          string
}

// BucketAPI interface for mocking
type BucketAPI interface {
	PutObject(objectKey string, reader io.Reader, options ...oss.Option) error
	DeleteObject(objectKey string, options ...oss.Option) error
	IsObjectExist(objectKey string, options ...oss.Option) (bool, error)
	ListObjects(options ...oss.Option) (oss.ListObjectsResult, error)
}

// realBucket wrapper
type realBucket struct {
	*oss.Bucket
}

// AliyunOSS implements Uploader for generic Alibaba Cloud OSS
type AliyunOSS struct {
	Client   *oss.Client
	Bucket   BucketAPI
	
	Cfg      OSSConfig
}

// NewAliyunOSS creates a new OSS uploader
// Precedence: config argument > environment variables -> error
func NewAliyunOSS(cfg *OSSConfig) (*AliyunOSS, error) {
	if cfg == nil {
		cfg = &OSSConfig{}
	}
	
	// Fallback to Env
	if cfg.Endpoint == "" {
		cfg.Endpoint = os.Getenv("ALIYUN_OSS_ENDPOINT")
	}
	if cfg.AccessKeyID == "" {
		cfg.AccessKeyID = os.Getenv("ALIYUN_ACCESS_KEY")
	}
	if cfg.AccessKeySecret == "" {
		cfg.AccessKeySecret = os.Getenv("ALIYUN_ACCESS_SECRET")
	}
	if cfg.BucketName == "" {
		cfg.BucketName = os.Getenv("ALIYUN_OSS_BUCKET")
	}
	if cfg.Domain == "" {
		cfg.Domain = os.Getenv("ALIYUN_OSS_DOMAIN")
	}

	// Validation
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.BucketName == "" {
		return nil, errors.New("aliyun oss config missing")
	}

	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to init oss client: %w", err)
	}

	bucket, err := client.Bucket(cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return &AliyunOSS{
		Client:   client,
		Bucket:   &realBucket{Bucket: bucket}, // Use wrapper
		Cfg:      *cfg,
	}, nil
}

// CreateFolder creates a directory marker (ends with /)
func (s *AliyunOSS) CreateFolder(name string) error {
	if !strings.HasSuffix(name, "/") {
		name += "/"
	}
	return s.Bucket.PutObject(name, bytes.NewReader([]byte("")))
}

func (s *AliyunOSS) IsExists(key string) (bool, error) {
	return s.Bucket.IsObjectExist(key)
}

func (s *AliyunOSS) Upload(key string, data io.Reader) (string, error) {
	// Simple Upload
	if err := s.Bucket.PutObject(key, data); err != nil {
		return "", fmt.Errorf("oss put failed: %w", err)
	}
	return s.GetURL(key), nil
}

func (s *AliyunOSS) Delete(key string) error {
	return s.Bucket.DeleteObject(key)
}

func (s *AliyunOSS) GetURL(key string) string {
	if s.Cfg.Domain != "" {
		// Use custom domain/CDN
		return strings.TrimRight(s.Cfg.Domain, "/") + "/" + strings.TrimLeft(key, "/")
	}
	// Default URL format: https://<bucket>.<endpoint>/<key>
	// Endpoint might contain http/https prefix or not
	endpoint := s.Cfg.Endpoint
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}
	// Note: Standard endpoint style is bucket.endpoint if not IP
	// This is a simplified estimation
	return fmt.Sprintf("https://%s.%s/%s", s.Cfg.BucketName, strings.TrimPrefix(endpoint, "https://"), key)
}
