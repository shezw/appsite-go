// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package simpleimage_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"appsite-go/pkg/utils/simpleimage"
)

func createTestImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with blue
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func TestValidateImage(t *testing.T) {
	imgData := createTestImage()
	
	// Valid image
	cfg, format, err := simpleimage.ValidateImage(bytes.NewReader(imgData))
	if err != nil {
		t.Fatalf("ValidateImage failed: %v", err)
	}
	if format != "png" {
		t.Errorf("Format = %s, want png", format)
	}
	if cfg.Width != 100 || cfg.Height != 100 {
		t.Errorf("Dims = %dx%d, want 100x100", cfg.Width, cfg.Height)
	}

	// Invalid image
	_, _, err = simpleimage.ValidateImage(bytes.NewReader([]byte("not an image")))
	if err == nil {
		t.Error("ValidateImage should fail for text data")
	}
}

func TestResizeImage(t *testing.T) {
	imgData := createTestImage()

	// Resize to 50x50
	resized, err := simpleimage.ResizeImage(bytes.NewReader(imgData), 50, 50)
	if err != nil {
		t.Fatalf("ResizeImage failed: %v", err)
	}
	if resized.Bounds().Dx() != 50 || resized.Bounds().Dy() != 50 {
		t.Errorf("Resized bounds = %v, want 50x50", resized.Bounds())
	}

	// Resize with one dim 0 (aspect ratio)
	resizedAspect, err := simpleimage.ResizeImage(bytes.NewReader(imgData), 50, 0)
	if err != nil {
		t.Fatalf("ResizeImage failed: %v", err)
	}
	if resizedAspect.Bounds().Dx() != 50 || resizedAspect.Bounds().Dy() != 50 {
		// Since original is 100x100, 50x0 should result in 50x50
		t.Errorf("Resized aspect bounds = %v, want 50x50", resizedAspect.Bounds())
	}
}

func TestThumbnail(t *testing.T) {
	imgData := createTestImage()

	// Thumbnail 20x20
	thumb, err := simpleimage.Thumbnail(bytes.NewReader(imgData), 20, 20)
	if err != nil {
		t.Fatalf("Thumbnail failed: %v", err)
	}
	if thumb.Bounds().Dx() != 20 || thumb.Bounds().Dy() != 20 {
		t.Errorf("Thumbnail bounds = %v, want 20x20", thumb.Bounds())
	}
}

func TestEncode(t *testing.T) {
	imgData := createTestImage()
	img, _, _ := image.Decode(bytes.NewReader(imgData))

	var buf bytes.Buffer
	err := simpleimage.Encode(&buf, img, "jpeg")
	if err != nil {
		t.Errorf("Encode jpeg failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("Encoded jpeg buffer is empty")
	}
}

func TestIsSupportedFormat(t *testing.T) {
	if !simpleimage.IsSupportedFormat("png") {
		t.Error("png should be supported")
	}
	if !simpleimage.IsSupportedFormat("JPG") {
		t.Error("JPG should be supported")
	}
	if simpleimage.IsSupportedFormat("bmp") {
		// In our check.go we only listed jpeg, jpg, png, gif
		t.Error("bmp should not be supported in strict check")
	}
}
