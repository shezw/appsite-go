// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package simpleimage

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

// ResizeImage resizes the image from reader to specified width and height
// If width or height is 0, it maintains aspect ratio
func ResizeImage(r io.Reader, width, height uint) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	// Lanczos3 provides high quality resizing
	m := resize.Resize(width, height, img, resize.Lanczos3)
	return m, nil
}

// Thumbnail creates a thumbnail of max width/height preserving aspect ratio
func Thumbnail(r io.Reader, maxWidth, maxHeight uint) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	m := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)
	return m, nil
}

// Encode saves the image to writer in specified format
func Encode(w io.Writer, img image.Image, format string) error {
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
	case "png":
		return png.Encode(w, img)
	case "gif":
		return gif.Encode(w, img, nil)
	default:
		// Default to JPEG if unknown, or maybe error out?
		// For simplicity, let's error or default. Choosing JPEG default for now.
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
	}
}
