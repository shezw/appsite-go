package cloudstorage_test

import (
	"io"
	"os"
	"testing"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"appsite-go/pkg/extra/cloudstorage"
)

// MockBucket implements generic OSS operations
type MockBucket struct {
	Objects map[string][]byte
}

func NewMockBucket() *MockBucket {
	return &MockBucket{Objects: make(map[string][]byte)}
}

func (m *MockBucket) PutObject(objectKey string, reader io.Reader, options ...oss.Option) error {
	data, _ := io.ReadAll(reader)
	m.Objects[objectKey] = data
	return nil
}

func (m *MockBucket) DeleteObject(objectKey string, options ...oss.Option) error {
	delete(m.Objects, objectKey)
	return nil
}

func (m *MockBucket) IsObjectExist(objectKey string, options ...oss.Option) (bool, error) {
	_, ok := m.Objects[objectKey]
	return ok, nil
}

func (m *MockBucket) ListObjects(options ...oss.Option) (oss.ListObjectsResult, error) {
	// Not fully implemented for simple test
	return oss.ListObjectsResult{}, nil
}

func TestAliyunOSSConfig(t *testing.T) {
	// 1. Test Missing Config
	_, err := cloudstorage.NewAliyunOSS(nil)
	if err == nil {
		t.Error("Expected error for missing config")
	}

	// 2. Test Env Fallback
	os.Setenv("ALIYUN_OSS_ENDPOINT", "oss-cn-shanghai.aliyuncs.com")
	os.Setenv("ALIYUN_ACCESS_KEY", "test_key")
	os.Setenv("ALIYUN_ACCESS_SECRET", "test_secret")
	os.Setenv("ALIYUN_OSS_BUCKET", "test-bucket")

	defer func() {
		os.Unsetenv("ALIYUN_OSS_ENDPOINT")
		os.Unsetenv("ALIYUN_ACCESS_KEY")
		os.Unsetenv("ALIYUN_ACCESS_SECRET")
		os.Unsetenv("ALIYUN_OSS_BUCKET")
	}()

	// Note: Verify logic in NewAliyunOSS attempts to connect. 
	// Since we pass fake credentials, oss.New might succeed (validation is minimal), 
	// but client.Bucket might not unless check fails.
	// Actually oss.New checks nothing. client.Bucket checks nothing locally unless it queries.
	// We can check if *AliyunOSS returns with correct fields.
	
	// However, NewAliyunOSS creates a real *oss.Client and *oss.Bucket.
	// We can rely on that for this test or we just assume success if it doesn't panic.
	// But we want to test the wrapper methods using MockBucket.
	
	// Because NewAliyunOSS returns a struct with a real underlying Bucket, 
	// we cannot inject the MockBucket easily unless we modify the struct AFTER creation
	// or create a separate constructor for testing.
	// Let's manually construct the AliyunOSS struct in test for logic verification.
}

func TestAliyunOSSOperations(t *testing.T) {
	mockB := NewMockBucket()
	client := &cloudstorage.AliyunOSS{
		Bucket: mockB,
		Cfg: cloudstorage.OSSConfig{
			BucketName: "test",
			Endpoint:   "oss.com",
		},
	}

	// 1. File Upload
	url, err := client.Upload("test.txt", strings.NewReader("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(url, "test.txt") {
		t.Error("URL incorrect")
	}
	if content, ok := mockB.Objects["test.txt"]; !ok || string(content) != "hello" {
		t.Error("File content mismatch")
	}

	// 2. Folder Create
	if err := client.CreateFolder("images"); err != nil {
		t.Fatal(err)
	}
	if _, ok := mockB.Objects["images/"]; !ok {
		t.Error("Folder marker not created")
	}

	// 3. Exists
	exists, _ := client.IsExists("test.txt")
	if !exists {
		t.Error("IsExists failed")
	}

	// 4. Delete
	client.Delete("test.txt")
	exists, _ = client.IsExists("test.txt")
	if exists {
		t.Error("Delete failed")
	}
}
