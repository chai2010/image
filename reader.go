// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"io"
	"os"
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

func LoadConfig(filename string) (cfg image.Config, format string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return image.Config{}, "", err
	}
	defer f.Close()
	return DecodeConfig(f)
}

func Load(filename string) (m image.Image, format string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return Decode(f)
}

func LoadImage(filename string) (m *Image, format string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	x, format, err := Decode(f)
	if err != nil {
		return nil, "", err
	}

	if m, _ = AsMemPImage(x); m == nil {
		m = NewImageFrom(x)
	}
	return m, format, nil
}
