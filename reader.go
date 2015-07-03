// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"io"
)

func DecodeConfig(r io.Reader) (cfg image.Config, format string, err error) {
	return image.DecodeConfig(r)
}

func Decode(r io.Reader) (m image.Image, format string, err error) {
	return image.Decode(r)
}

func DecodeImage(r io.Reader) (m *Image, format string, err error) {
	x, format, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}
	m = NewImageFrom(x)
	return
}

func Load(filename string) (m image.Image, format string, err error) {
	panic("TODO")
}

func LoadImage(filename string) (m *Image, format string, err error) {
	panic("TODO")
}
