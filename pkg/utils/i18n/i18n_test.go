package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestI18n(t *testing.T) {
	// Setup temporary config dir
	tmpDir, err := os.MkdirTemp("", "i18n_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create en-US.yaml
	enUS := []byte(`
common:
  hello: "Hello"
  welcome: "Welcome, %s"
`)
	if err := os.WriteFile(filepath.Join(tmpDir, "en-US.yaml"), enUS, 0644); err != nil {
		t.Fatal(err)
	}

	// Create zh-CN.yaml
	zhCN := []byte(`
common:
  hello: "你好"
`)
	if err := os.WriteFile(filepath.Join(tmpDir, "zh-CN.yaml"), zhCN, 0644); err != nil {
		t.Fatal(err)
	}

	// Init
	if err := Init(tmpDir, "en-US"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		lang     string
		key      string
		args     []interface{}
		expected string
	}{
		{"en-US", "common.hello", nil, "Hello"},
		{"zh-CN", "common.hello", nil, "你好"},
		{"fr-FR", "common.hello", nil, "Hello"}, // Fallback to en-US
		{"en-US", "common.welcome", []interface{}{"John"}, "Welcome, John"}, 
		{"en-US", "common.unknown", nil, "common.unknown"},
	}

	for _, tt := range tests {
		got := T(tt.lang, tt.key, tt.args...)
		if got != tt.expected {
			t.Errorf("T(%q, %q) = %q, want %q", tt.lang, tt.key, got, tt.expected)
		}
	}
}
