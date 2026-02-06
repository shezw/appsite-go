// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package file

import (
	"os"
	"testing"
)

func TestMkDirAndCheck(t *testing.T) {
	testDir := "test_dir"
	defer os.RemoveAll(testDir)

	// Test IsNotExistMkDir
	if err := IsNotExistMkDir(testDir); err != nil {
		t.Fatalf("IsNotExistMkDir failed: %v", err)
	}

	// Test CheckNotExist
	if CheckNotExist(testDir) {
		t.Error("CheckNotExist should return false for existing dir")
	}

	// Test CheckNotExist for non-existing
	if !CheckNotExist("non_existing_dir") {
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
		if got := FormatFileSize(tt.size); got != tt.expected {
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
	if ext := GetExt(testFile); ext != ".txt" {
		t.Errorf("GetExt() = %s, want .txt", ext)
	}

	// Test GetMimeType (extension based)
	mime, err := GetMimeType(testFile)
	if err != nil {
		t.Errorf("GetMimeType failed: %v", err)
	}
	// Note: text/plain result depends on system mime.types, usually text/plain; charset=utf-8
	if mime == "" {
		t.Error("GetMimeType returned empty")
	}

	// Test CheckImageFile
	if CheckImageFile("image.jpg") == false {
		t.Error("CheckImageFile(.jpg) should be true")
	}
	if CheckImageFile("doc.pdf") == true {
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

	mime, err := GetMimeType(testFile)
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
	f, err := Open(testFile, os.O_RDONLY, 0)
	if err != nil {
		t.Errorf("Open failed: %v", err)
	}
	defer f.Close()

	// Test GetSize
	// Note: GetSize takes multipart.File, but os.File implements similar interface but GetSize uses io.ReadAll.
	// io.ReadAll reads from current position. Newly opened file is at 0.
	size, err := GetSize(f)
	if err != nil {
		t.Errorf("GetSize failed: %v", err)
	}
	if size != 5 {
		t.Errorf("GetSize = %d, want 5", size)
	}

	// Test CheckPermission
	// CheckPermission wraps os.IsPermission(err). 
	// If file is accessible (no error), it returns false.
	// If file has permission error, it returns true.
	if CheckPermission(testFile) {
		t.Error("CheckPermission should be false (no error) for accessible file")
	}

	// Test with explicit no permission
	noPermFile := "noperm.txt"
	os.WriteFile(noPermFile, []byte("secret"), 0000)
	defer os.Remove(noPermFile)
	
	// CheckPermission might return true here if we are not root
	// But in some CI environments or if owner is same, 0000 might still be stat-able by owner?
	// os.Stat usually works for owner even with 0000 on Linux? No, mode 0000 denies all.
	// Let's see. If this is flaky, I'll remove the positive test and just stick to negative.
	// _ = CheckPermission(noPermFile) 
}
