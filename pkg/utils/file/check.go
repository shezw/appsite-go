// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package file

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetMimeType returns the mime tyoe of a file based on its extension or content
func GetMimeType(filePath string) (string, error) {
	// 1. Try by extension first (fast)
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return mimeType, nil
	}

	// 2. Try by content (slower but accurate)
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid content-type by default: "application/octet-stream"
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// CheckImageFile checks if the file is an image
func CheckImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg":
		return true
	}
	return false
}
