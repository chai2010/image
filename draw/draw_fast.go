// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"
	"image/draw"

	ximage "github.com/chai2010/image"
	xdraw "golang.org/x/image/draw"
)

func drawFast(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	switch dst := dst.(type) {
	case *image.Gray:
		if src, ok := src.(*image.Gray); ok {
			drawGray(dst, r, src, sp)
			return
		}
	case *image.Gray16:
		if src, ok := src.(*image.Gray16); ok {
			drawGray16(dst, r, src, sp)
			return
		}
	case *image.RGBA:
		if src, ok := src.(*image.RGBA); ok {
			drawRGBA(dst, r, src, sp)
			return
		}
	case *image.RGBA64:
		if src, ok := src.(*image.RGBA64); ok {
			drawRGBA64(dst, r, src, sp)
			return
		}
	}

	if dst, ok := ximage.AsMemPImage(dst); ok {
		if src, ok := ximage.AsMemPImage(src); ok {
			drawImage(dst, r, src, sp)
			return
		}
	}

	xdraw.Draw(dst, r, src, sp, xdraw.Src)
}

func drawGray(dst *image.Gray, r image.Rectangle, src *image.Gray, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*1*1]
		srcLine := src.Pix[off1:][:r.Dx()*1*1]

		copy(dstLine, srcLine)
	}
}

func drawGray16(dst *image.Gray16, r image.Rectangle, src *image.Gray16, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*2*1]
		srcLine := src.Pix[off1:][:r.Dx()*2*1]

		copy(dstLine, srcLine)
	}
}

func drawRGBA(dst *image.RGBA, r image.Rectangle, src *image.RGBA, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*1*4]
		srcLine := src.Pix[off1:][:r.Dx()*1*4]

		copy(dstLine, srcLine)
	}
}

func drawRGBA64(dst *image.RGBA64, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*2*4]
		srcLine := src.Pix[off1:][:r.Dx()*2*4]

		copy(dstLine, srcLine)
	}
}

func drawImage(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	if dst.Channels != src.Channels || dst.DataType != src.DataType {
		xdraw.Draw(dst, r, src, sp, xdraw.Src)
		return
	}

	dx := ximage.SizeofPixel(dst.Channels, dst.DataType) * r.Dx()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:dx]
		srcLine := src.Pix[off1:][:dx]

		copy(dstLine, srcLine)
	}
}
