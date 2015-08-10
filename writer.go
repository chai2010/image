// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Encoder func(w io.Writer, m image.Image) error

// Encode image, if encoder is nil, use png format.
func Encode(m image.Image, encoder Encoder) ([]byte, error) {
	if encoder == nil {
		encoder = png.Encode
	}
	b := bytes.NewBuffer([]byte{})
	if err := encoder(b, m); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Save image, if encoder is nil, only support gif/jpeg/png format.
func Save(filename string, m image.Image, encoder Encoder) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if encoder != nil {
		return encoder(f, m)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	switch {
	case strings.HasSuffix(ext, ".gif"):
		return gif.Encode(f, m, nil)
	case strings.HasSuffix(ext, ".jpeg"):
		return jpeg.Encode(f, m, nil)
	case strings.HasSuffix(ext, ".jpg"):
		return jpeg.Encode(f, m, nil)
	case strings.HasSuffix(ext, ".png"):
		return png.Encode(f, m)
	}

	return fmt.Errorf("image: Save, unknown format: %s", filename)
}
