// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package simpleimage

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	kerror "appsite-go/internal/core/error"
)

// ValidateImage checks if the reader contains a valid image and returns its config
func ValidateImage(r io.Reader) (image.Config, string, error) {
	config, format, err := image.DecodeConfig(r)
	if err != nil {
		return image.Config{}, "", kerror.NewWithMessage(kerror.InvalidParams, "invalid image format")
	}
	return config, format, nil
}

// IsSupportedFormat checks if the format string is one of the supported ones
func IsSupportedFormat(format string) bool {
	switch strings.ToLower(format) {
	case "jpeg", "jpg", "png", "gif":
		return true
	}
	return false
}
