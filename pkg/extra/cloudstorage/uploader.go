// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cloudstorage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Uploader defines the interface for file storage
type Uploader interface {
	// Upload saves data to the storage backend and returns the public URL
	Upload(key string, data io.Reader) (url string, err error)
	// Delete removes the file
	Delete(key string) error
	// GetURL returns the public access URL for a key
	GetURL(key string) string
}

// LocalStorage implements Uploader for the local filesystem
type LocalStorage struct {
	BaseDir   string // Local directory path, e.g. "./uploads"
	BaseURL   string // Public URL prefix, e.g. "http://localhost:8080/uploads"
}

func NewLocalStorage(baseDir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{
		BaseDir: baseDir,
		BaseURL: baseURL,
	}, nil
}

func (s *LocalStorage) Upload(key string, data io.Reader) (string, error) {
	// Security check: simple path traversal prevention
	// real implementation might need more robust checks
	cleanKey := filepath.Clean(key) 
	dstPath := filepath.Join(s.BaseDir, cleanKey)
	
	// Ensure directory exists for nested keys
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return "", err
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, data)
	if err != nil {
		return "", err
	}

	return s.BaseURL + "/" + cleanKey, nil
}

func (s *LocalStorage) Delete(key string) error {
	cleanKey := filepath.Clean(key)
	dstPath := filepath.Join(s.BaseDir, cleanKey)
	return os.Remove(dstPath)
}

func (s *LocalStorage) GetURL(key string) string {
	return s.BaseURL + "/" + filepath.Clean(key)
}

type MockS3Storage struct {
	Bucket string
	Region string
}

func (s *MockS3Storage) Upload(key string, data io.Reader) (string, error) {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.Bucket, s.Region, key), nil
}

func (s *MockS3Storage) Delete(key string) error {
	return nil
}

func (s *MockS3Storage) GetURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.Bucket, s.Region, key)
}
