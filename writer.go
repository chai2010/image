// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
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
