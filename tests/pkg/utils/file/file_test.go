// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package file_test

import (
	"os"
	"testing"

	"appsite-go/pkg/utils/file"
)

func TestMkDirAndCheck(t *testing.T) {
	testDir := "test_dir"
	defer os.RemoveAll(testDir)

	// Test IsNotExistMkDir
	if err := file.IsNotExistMkDir(testDir); err != nil {
		t.Fatalf("IsNotExistMkDir failed: %v", err)
	}

	// Test CheckNotExist
	if file.CheckNotExist(testDir) {
		t.Error("CheckNotExist should return false for existing dir")
	}

	// Test CheckNotExist for non-existing
	if !file.CheckNotExist("non_existing_dir") {
		t.Error("CheckNotExist should return true for non-existing dir")
	}
}

func TestFileSizeFormat(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{100, "100.00 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
	}

	for _, tt := range tests {
		if got := file.FormatFileSize(tt.size); got != tt.expected {
			t.Errorf("FormatFileSize(%d) = %s, want %s", tt.size, got, tt.expected)
		}
	}
}

func TestMimeAndExt(t *testing.T) {
	// Create a dummy file
	testFile := "test.txt"
	content := []byte("hello world")
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	// Test GetExt
	if ext := file.GetExt(testFile); ext != ".txt" {
		t.Errorf("GetExt() = %s, want .txt", ext)
	}

	// Test GetMimeType (extension based)
	mime, err := file.GetMimeType(testFile)
	if err != nil {
		t.Errorf("GetMimeType failed: %v", err)
	}
	if mime == "" {
		t.Error("GetMimeType returned empty")
	}

	// Test CheckImageFile
	if file.CheckImageFile("image.jpg") == false {
		t.Error("CheckImageFile(.jpg) should be true")
	}
	if file.CheckImageFile("doc.pdf") == true {
		t.Error("CheckImageFile(.pdf) should be false")
	}
}

func TestGetMimeTypeContent(t *testing.T) {
	// Create a file without extension but png content
	testFile := "test_png"
	// Small PNG header
	content := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	mime, err := file.GetMimeType(testFile)
	if err != nil {
		t.Errorf("GetMimeType failed: %v", err)
	}
	if mime != "image/png" {
		t.Errorf("GetMimeType = %s, want image/png", mime)
	}
}

func TestFileHelpers(t *testing.T) {
	testFile := "helper_test.txt"
	content := []byte("12345")
	os.WriteFile(testFile, content, 0644)
	defer os.Remove(testFile)

	// Test Open
	f, err := file.Open(testFile, os.O_RDONLY, 0)
	if err != nil {
		t.Errorf("Open failed: %v", err)
	}
	defer f.Close()

	// Test GetSize
	size, err := file.GetSize(f)
	if err != nil {
		t.Errorf("GetSize failed: %v", err)
	}
	if size != 5 {
		t.Errorf("GetSize = %d, want 5", size)
	}

	// Test CheckPermission
	if file.CheckPermission(testFile) {
		t.Error("CheckPermission should be false (no error) for accessible file")
	}
}
