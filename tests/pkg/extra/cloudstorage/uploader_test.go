package cloudstorage_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"appsite-go/pkg/extra/cloudstorage"
)

func TestLocalStorage(t *testing.T) {
	tmpDir := "./tmp_uploads"
	defer os.RemoveAll(tmpDir)

	storage, err := cloudstorage.NewLocalStorage(tmpDir, "http://cdn.example.com")
	if err != nil {
		t.Fatal(err)
	}

	key := "images/avatar.png"
	content := []byte("fake png content")
	
	// Test Upload
	url, err := storage.Upload(key, bytes.NewReader(content))
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	expectedPrefix := "http://cdn.example.com/"
	if !strings.HasPrefix(url, expectedPrefix) {
		t.Errorf("URL format wrong: %s", url)
	}

	// Verify file exists
	if _, err := os.Stat(tmpDir + "/" + key); os.IsNotExist(err) {
		t.Error("File not created on disk")
	}

	// Test GetURL
	if u := storage.GetURL(key); u != url {
		t.Errorf("GetURL mismatch: %s vs %s", u, url)
	}

	// Test Delete
	if err := storage.Delete(key); err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	
	if _, err := os.Stat(tmpDir + "/" + key); !os.IsNotExist(err) {
		t.Error("File still exists after delete")
	}
}

func TestMockS3(t *testing.T) {
	s3 := &cloudstorage.MockS3Storage{
		Bucket: "mybucket",
		Region: "us-east-1",
	}
	
	url, err := s3.Upload("test.jpg", strings.NewReader("data"))
	if err != nil {
		t.Fatal(err)
	}
	
	if !strings.Contains(url, "mybucket.s3.us-east-1.amazonaws.com") {
		t.Errorf("URL invalid: %s", url)
	}
}
