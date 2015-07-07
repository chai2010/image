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

func nnPyrDownFast(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	switch dst := dst.(type) {
	case *image.Gray:
		if src, ok := src.(*image.Gray); ok {
			nnPyrDownGray(dst, r, src, sp)
			return
		}
	case *image.Gray16:
		if src, ok := src.(*image.Gray16); ok {
			nnPyrDownGray16(dst, r, src, sp)
			return
		}
	case *image.RGBA:
		if src, ok := src.(*image.RGBA); ok {
			nnPyrDownRGBA(dst, r, src, sp)
			return
		}
	case *image.RGBA64:
		if src, ok := src.(*image.RGBA64); ok {
			nnPyrDownRGBA64(dst, r, src, sp)
			return
		}
	}

	if dst, ok := ximage.AsMemPImage(dst); ok {
		if src, ok := ximage.AsMemPImage(src); ok {
			nnPyrDownImage(dst, r, src, sp)
			return
		}
	}

	xdraw.NearestNeighbor.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func nnPyrDownGray(dst *image.Gray, r image.Rectangle, src *image.Gray, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLine := dst.Pix[off0:][:r.Dx()*1*1]
		srcLine := src.Pix[off1:][:r.Dx()*1*2]

		for i, j := 0, 0; i < len(dstLine); i, j = i+1*1, j+1*2 {
			dstLine[i] = srcLine[j]
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
	}
}

func nnPyrDownGray16(dst *image.Gray16, r image.Rectangle, src *image.Gray16, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLine := dst.Pix[off0:][:r.Dx()*2*1]
		srcLine := src.Pix[off1:][:r.Dx()*2*2]

		for i, j := 0, 0; i < len(dstLine); i, j = i+2*1, j+2*2 {
			dstLine[i+0] = srcLine[j+0]
			dstLine[i+1] = srcLine[j+1]
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
	}
}

func nnPyrDownRGBA(dst *image.RGBA, r image.Rectangle, src *image.RGBA, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLine := dst.Pix[off0:][:r.Dx()*4*1]
		srcLine := src.Pix[off1:][:r.Dx()*4*2]

		for i, j := 0, 0; i < len(dstLine); i, j = i+4*1, j+4*2 {
			dstLine[i+0] = srcLine[j+0]
			dstLine[i+1] = srcLine[j+1]
			dstLine[i+2] = srcLine[j+2]
			dstLine[i+3] = srcLine[j+3]
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
	}
}

func nnPyrDownRGBA64(dst *image.RGBA64, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLine := dst.Pix[off0:][:r.Dx()*8*1]
		srcLine := src.Pix[off1:][:r.Dx()*8*2]

		for i, j := 0, 0; i < len(dstLine); i, j = i+8*1, j+8*2 {
			dstLine[i+0] = srcLine[j+0]
			dstLine[i+1] = srcLine[j+1]
			dstLine[i+2] = srcLine[j+2]
			dstLine[i+3] = srcLine[j+3]
			dstLine[i+4] = srcLine[j+4]
			dstLine[i+5] = srcLine[j+5]
			dstLine[i+6] = srcLine[j+6]
			dstLine[i+7] = srcLine[j+7]
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
	}
}

func nnPyrDownImage(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	if dst.Channels != src.Channels || dst.DataType != src.DataType {
		xdraw.NearestNeighbor.Scale(
			dst, r,
			src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
			xdraw.Src, nil,
		)
		return
	}

	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	pixSize := ximage.SizeofPixel(dst.Channels, dst.DataType)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLine := dst.Pix[off0:][:r.Dx()*pixSize*1]
		srcLine := src.Pix[off1:][:r.Dx()*pixSize*2]

		for i, j := 0, 0; i < len(dstLine); i, j = i+pixSize*1, j+pixSize*2 {
			copy(dstLine[i:i+pixSize], srcLine[j:j+pixSize])
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
	}
}
