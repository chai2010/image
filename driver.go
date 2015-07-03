// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"io"
)

type Driver interface {
	Name() string
	Formats() []string
	Extensions() []string

	MatchHeader(header []byte) (strongMatched, weakMatched bool)
	DecodeConfig(r io.Reader) (cfg image.Config, format string, err error)
	Decode(r io.Reader) (m image.Image, format string, err error)
	Encode(w io.Writer, m image.Image, format string) error

	MatchFile(filename string) (strongMatched, weakMatched bool)
	LoadHeader(filename string) (cfg image.Config, format string, err error)
	Load(filename string) (m image.Image, format string, err error)
	Save(filename string, m image.Image, format string) error
}

func RegisterDriver(driver Driver) {
	panic("TODO")
}
