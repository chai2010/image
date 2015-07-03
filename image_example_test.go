// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image_test

import (
	"image"
	"reflect"

	memp "."
)

func Example_rgb() {
	type RGB struct {
		R, G, B uint8
	}

	b := image.Rect(0, 0, 300, 400)
	rgbImage := memp.NewImage(b, 3, reflect.Uint8)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		var (
			line     []byte = rgbImage.Pix[rgbImage.PixOffset(b.Min.X, y):][:rgbImage.Stride]
			rgbSlice []RGB  = memp.PixSilce(line).Slice(reflect.TypeOf([]RGB(nil))).([]RGB)
		)

		for i, _ := range rgbSlice {
			rgbSlice[i] = RGB{
				R: uint8(i + 1),
				G: uint8(i + 2),
				B: uint8(i + 3),
			}
		}
	}
}

func Example_rgb48() {
	type RGB struct {
		R, G, B uint16
	}

	b := image.Rect(0, 0, 300, 400)
	rgbImage := memp.NewImage(b, 3, reflect.Uint16)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		var (
			line     []byte = rgbImage.Pix[rgbImage.PixOffset(b.Min.X, y):][:rgbImage.Stride]
			rgbSlice []RGB  = memp.PixSilce(line).Slice(reflect.TypeOf([]RGB(nil))).([]RGB)
		)

		for i, _ := range rgbSlice {
			rgbSlice[i] = RGB{
				R: uint16(i + 1),
				G: uint16(i + 2),
				B: uint16(i + 3),
			}
		}
	}
}

func Example_unsafe() {
	// struct must same as memp.Image
	type MyImage struct {
		MemPMagic string // MemP
		Rect      image.Rectangle
		Channels  int
		DataType  reflect.Kind
		Pix       []byte

		// Stride is the Pix stride (in bytes, must align with PixelSize)
		// between vertically adjacent pixels.
		Stride int
	}

	p := &MyImage{
		MemPMagic: memp.MemPMagic,
		Rect:      image.Rect(0, 0, 300, 400),
		Channels:  3,
		DataType:  reflect.Uint16,
		Pix:       make([]byte, 300*400*3*memp.SizeofKind(reflect.Uint16)),
		Stride:    300 * 3 * memp.SizeofKind(reflect.Uint16),
	}

	q, _ := memp.AsMemPImage(p)
	_ = q
}
