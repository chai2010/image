// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"io"
	"os"
)

type LoadConfiger func(filename string) (cfg image.Config, format string, err error)
type Loader func(filename string) (m image.Image, format string, err error)

func DecodeConfig(r io.Reader) (cfg image.Config, format string, err error) {
	return image.DecodeConfig(r)
}

func Decode(r io.Reader) (m image.Image, format string, err error) {
	return image.Decode(r)
}

func DecodeImage(r io.Reader) (m *MemPImage, format string, err error) {
	x, format, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}
	m = NewMemPImageFrom(x)
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

func LoadConfigEx(filename string, loader LoadConfiger) (cfg image.Config, format string, err error) {
	if loader != nil {
		cfg, format, err = loader(filename)
		if err != nil {
			var err1 error
			if cfg, format, err1 = LoadConfig(filename); err1 != nil {
				return image.Config{}, "", err // retuen loader's err
			}
		}
		return
	}
	return LoadConfig(filename)
}

func Load(filename string) (m image.Image, format string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return Decode(f)
}

func LoadEx(filename string, loader Loader) (m image.Image, format string, err error) {
	if loader != nil {
		m, format, err = loader(filename)
		if err != nil {
			var err1 error
			if m, format, err1 = Load(filename); err1 != nil {
				return nil, "", err // retuen loader's err
			}
		}
		return
	}
	return Load(filename)
}

func LoadImage(filename string) (m *MemPImage, format string, err error) {
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
		m = NewMemPImageFrom(x)
	}
	return m, format, nil
}

func LoadImageEx(filename string, loader Loader) (m *MemPImage, format string, err error) {
	if loader != nil {
		x, format, err := loader(filename)
		if err != nil {
			return nil, "", err
		}
		if m, _ = AsMemPImage(x); m == nil {
			m = NewMemPImageFrom(x)
		}
		return m, format, nil
	}
	return LoadImage(filename)
}
