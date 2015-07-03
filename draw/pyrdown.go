// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"
	"image/color"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

var (
	// NearestNeighbor is the nearest neighbor interpolator. It is very fast,
	// but usually gives very low quality results. When scaling up, the result
	// will look 'blocky'.
	NearestNeighbor = PyrDowner(nnPyrDowner{})

	// ApproxBiLinear is a mixture of the nearest neighbor and bi-linear
	// interpolators. It is fast, but usually gives medium quality results.
	//
	// It implements bi-linear interpolation when upscaling and a bi-linear
	// blend of the 4 nearest neighbor pixels when downscaling. This yields
	// nicer quality than nearest neighbor interpolation when upscaling, but
	// the time taken is independent of the number of source pixels, unlike the
	// bi-linear interpolator. When downscaling a large image, the performance
	// difference can be significant.
	ApproxBiLinear = PyrDowner(abPyrDowner{})
)

type PyrDowner interface {
	PyrDown(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point)
}

func Draw(dst xdraw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	// TODO
}

type nnPyrDowner struct{}

func (nnPyrDowner) PyrDown(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	xdraw.NearestNeighbor.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()*2, sp.Y+r.Dy()*2),
		xdraw.Src, nil,
	)
}

type abPyrDowner struct{}

func (abPyrDowner) PyrDown(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	_PyrDown_ApproxBiLinear(dst, r, src, sp)
}

func _PyrDown_ApproxBiLinear_Gray_Gray(dst *image.Gray, r image.Rectangle, src *image.Gray, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*1*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*1*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*1*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+1*1, j+1*2 {
			v00 := srcLine0[j+0]
			v01 := srcLine0[j+1]
			v10 := srcLine1[j+0]
			v11 := srcLine1[j+1]
			dstLineX[i] = uint8(((uint32(v00) + uint32(v01) + uint32(v10) + uint32(v11)) * 0x101) >> 10)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func _PyrDown_ApproxBiLinear_Gray_Gray16(dst *image.Gray, r image.Rectangle, src *image.Gray16, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*1*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*2*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*2*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+1*1, j+2*2 {
			v00 := uint32(srcLine0[j+0])<<8 | uint32(srcLine0[j+1])
			v01 := uint32(srcLine0[j+2])<<8 | uint32(srcLine0[j+3])
			v10 := uint32(srcLine1[j+0])<<8 | uint32(srcLine1[j+1])
			v11 := uint32(srcLine1[j+2])<<8 | uint32(srcLine1[j+3])

			dstLineX[i] = uint8((uint32(v00) + uint32(v01) + uint32(v10) + uint32(v11)) >> 10)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func _PyrDown_ApproxBiLinear_Gray_RGBA(dst *image.Gray, r image.Rectangle, src *image.RGBA, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*1*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*4*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*4*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+1*1, j+4*2 {
			r00 := srcLine0[j+0]
			g00 := srcLine0[j+1]
			b00 := srcLine0[j+2]

			r01 := srcLine0[j+4]
			g01 := srcLine0[j+5]
			b01 := srcLine0[j+6]

			r10 := srcLine1[j+0]
			g10 := srcLine1[j+1]
			b10 := srcLine1[j+2]

			r11 := srcLine1[j+4]
			g11 := srcLine1[j+5]
			b11 := srcLine1[j+6]

			rxx := (uint32(r00) + uint32(r01) + uint32(r10) + uint32(r11)) * 0x101
			gxx := (uint32(g00) + uint32(g01) + uint32(g10) + uint32(g11)) * 0x101
			bxx := (uint32(b00) + uint32(b01) + uint32(b10) + uint32(b11)) * 0x101

			dstLineX[i] = uint8(((299*rxx + 587*gxx + 114*bxx + 500) / 1000) >> 8)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func _PyrDown_ApproxBiLinear_Gray_RGBA64(dst *image.Gray, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_Gray_YCbCr(dst *image.Gray, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	_PyrDown_ApproxBiLinear_Gray_Gray(dst, r, &image.Gray{
		Pix:    src.Y,
		Stride: src.YStride,
		Rect:   src.Rect,
	}, sp)
}

func _PyrDown_ApproxBiLinear_Gray16_Gray(dst *image.Gray16, r image.Rectangle, src *image.Gray, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_Gray16_Gray16(dst *image.Gray16, r image.Rectangle, src *image.Gray16, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*2*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*2*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*2*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+2*1, j+2*2 {
			v00 := uint32(srcLine0[j+0])<<8 | uint32(srcLine0[j+1])
			v01 := uint32(srcLine0[j+2])<<8 | uint32(srcLine0[j+3])
			v10 := uint32(srcLine1[j+0])<<8 | uint32(srcLine1[j+1])
			v11 := uint32(srcLine1[j+2])<<8 | uint32(srcLine1[j+3])

			vxx := uint16((v00 + v01 + v10 + v11) / 4)
			dstLineX[i+0] = uint8(vxx >> 8)
			dstLineX[i+1] = uint8(vxx)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func _PyrDown_ApproxBiLinear_Gray16_RGBA(dst *image.Gray16, r image.Rectangle, src *image.RGBA, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_Gray16_RGBA64(dst *image.Gray16, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_Gray16_YCbCr(dst *image.Gray16, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_RGBA_Gray(dst *image.RGBA, r image.Rectangle, src *image.Gray, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_RGBA_Gray16(dst *image.RGBA, r image.Rectangle, src *image.Gray16, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_RGBA_RGBA(dst *image.RGBA, r image.Rectangle, src *image.RGBA, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*4*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*4*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*4*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+4*1, j+4*2 {
			r00 := srcLine0[j+0]
			g00 := srcLine0[j+1]
			b00 := srcLine0[j+2]
			a00 := srcLine0[j+3]

			r01 := srcLine0[j+4]
			g01 := srcLine0[j+5]
			b01 := srcLine0[j+6]
			a01 := srcLine0[j+7]

			r10 := srcLine1[j+0]
			g10 := srcLine1[j+1]
			b10 := srcLine1[j+2]
			a10 := srcLine1[j+3]

			r11 := srcLine1[j+4]
			g11 := srcLine1[j+5]
			b11 := srcLine1[j+6]
			a11 := srcLine1[j+7]

			dstLineX[i+0] = uint8(((uint32(r00) + uint32(r01) + uint32(r10) + uint32(r11)) * 0x101) >> 10)
			dstLineX[i+1] = uint8(((uint32(g00) + uint32(g01) + uint32(g10) + uint32(g11)) * 0x101) >> 10)
			dstLineX[i+2] = uint8(((uint32(b00) + uint32(b01) + uint32(b10) + uint32(b11)) * 0x101) >> 10)
			dstLineX[i+3] = uint8(((uint32(a00) + uint32(a01) + uint32(a10) + uint32(a11)) * 0x101) >> 10)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func _PyrDown_ApproxBiLinear_RGBA_RGBA64(dst *image.RGBA, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_RGBA_YCbCr(dst *image.RGBA, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_RGBA64_Gray(dst *image.RGBA64, r image.Rectangle, src *image.Gray, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_RGBA64_Gray16(dst *image.RGBA64, r image.Rectangle, src *image.Gray16, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_RGBA64_RGBA(dst *image.RGBA64, r image.Rectangle, src *image.RGBA, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear_RGBA64_RGBA64(dst *image.RGBA64, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*8*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*8*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*8*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+8*1, j+8*2 {
			r00 := uint32(srcLine0[j+8*0+0])<<8 | uint32(srcLine0[j+8*0+1])
			g00 := uint32(srcLine0[j+8*0+2])<<8 | uint32(srcLine0[j+8*0+3])
			b00 := uint32(srcLine0[j+8*0+4])<<8 | uint32(srcLine0[j+8*0+5])
			a00 := uint32(srcLine0[j+8*0+6])<<8 | uint32(srcLine0[j+8*0+7])
			r01 := uint32(srcLine0[j+8*1+0])<<8 | uint32(srcLine0[j+8*1+1])
			g01 := uint32(srcLine0[j+8*1+2])<<8 | uint32(srcLine0[j+8*1+3])
			b01 := uint32(srcLine0[j+8*1+4])<<8 | uint32(srcLine0[j+8*1+5])
			a01 := uint32(srcLine0[j+8*1+6])<<8 | uint32(srcLine0[j+8*1+7])
			r10 := uint32(srcLine1[j+8*0+0])<<8 | uint32(srcLine1[j+8*0+1])
			g10 := uint32(srcLine1[j+8*0+2])<<8 | uint32(srcLine1[j+8*0+3])
			b10 := uint32(srcLine1[j+8*0+4])<<8 | uint32(srcLine1[j+8*0+5])
			a10 := uint32(srcLine1[j+8*0+6])<<8 | uint32(srcLine1[j+8*0+7])
			r11 := uint32(srcLine1[j+8*1+0])<<8 | uint32(srcLine1[j+8*1+1])
			g11 := uint32(srcLine1[j+8*1+2])<<8 | uint32(srcLine1[j+8*1+3])
			b11 := uint32(srcLine1[j+8*1+4])<<8 | uint32(srcLine1[j+8*1+5])
			a11 := uint32(srcLine1[j+8*1+6])<<8 | uint32(srcLine1[j+8*1+7])

			rxx := uint16((r00 + r01 + r10 + r11) / 4)
			gxx := uint16((g00 + g01 + g10 + g11) / 4)
			bxx := uint16((b00 + b01 + b10 + b11) / 4)
			axx := uint16((a00 + a01 + a10 + a11) / 4)

			dstLineX[i+0] = uint8(rxx >> 8)
			dstLineX[i+1] = uint8(rxx)
			dstLineX[i+2] = uint8(gxx >> 8)
			dstLineX[i+3] = uint8(gxx)
			dstLineX[i+4] = uint8(bxx >> 8)
			dstLineX[i+5] = uint8(bxx)
			dstLineX[i+6] = uint8(axx >> 8)
			dstLineX[i+7] = uint8(axx)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func _PyrDown_ApproxBiLinear_RGBA64_YCbCr(dst *image.RGBA64, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func _PyrDown_ApproxBiLinear(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	switch dst := dst.(type) {
	case *image.Gray:
		switch src := src.(type) {
		case *image.Gray:
			_PyrDown_ApproxBiLinear_Gray_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			_PyrDown_ApproxBiLinear_Gray_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			_PyrDown_ApproxBiLinear_Gray_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			_PyrDown_ApproxBiLinear_Gray_RGBA64(dst, r, src, sp)
			return
		case *image.YCbCr:
			_PyrDown_ApproxBiLinear_Gray_YCbCr(dst, r, src, sp)
			return
		}
	case *image.Gray16:
		switch src := src.(type) {
		case *image.Gray:
			_PyrDown_ApproxBiLinear_Gray16_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			_PyrDown_ApproxBiLinear_Gray16_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			_PyrDown_ApproxBiLinear_Gray16_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			_PyrDown_ApproxBiLinear_Gray16_RGBA64(dst, r, src, sp)
			return
		case *image.YCbCr:
			_PyrDown_ApproxBiLinear_Gray16_YCbCr(dst, r, src, sp)
			return
		}
	case *image.RGBA:
		switch src := src.(type) {
		case *image.Gray:
			_PyrDown_ApproxBiLinear_RGBA_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			_PyrDown_ApproxBiLinear_RGBA_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			_PyrDown_ApproxBiLinear_RGBA_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			_PyrDown_ApproxBiLinear_RGBA_RGBA64(dst, r, src, sp)
			return
		case *image.YCbCr:
			_PyrDown_ApproxBiLinear_RGBA_YCbCr(dst, r, src, sp)
			return
		}
	case *image.RGBA64:
		switch src := src.(type) {
		case *image.Gray:
			_PyrDown_ApproxBiLinear_RGBA64_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			_PyrDown_ApproxBiLinear_RGBA64_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			_PyrDown_ApproxBiLinear_RGBA64_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			_PyrDown_ApproxBiLinear_RGBA64_RGBA64(dst, r, src, sp)
			return
		case *image.YCbCr:
			_PyrDown_ApproxBiLinear_RGBA64_YCbCr(dst, r, src, sp)
			return
		}
	}

	dstColorRGBA64 := &color.RGBA64{}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			sx := sp.X + (x-r.Min.X)*2
			sy := sp.Y + (y-r.Min.Y)*2

			r00, g00, b00, a00 := src.At(sx+0, sy+0).RGBA()
			r01, g01, b01, a01 := src.At(sx+0, sy+1).RGBA()
			r11, g11, b11, a11 := src.At(sx+1, sy+1).RGBA()
			r10, g10, b10, a10 := src.At(sx+1, sy+0).RGBA()

			dstColorRGBA64.R = uint16((r00 + r01 + r11 + r10) / 4)
			dstColorRGBA64.G = uint16((g00 + g01 + g11 + g10) / 4)
			dstColorRGBA64.B = uint16((b00 + b01 + b11 + b10) / 4)
			dstColorRGBA64.A = uint16((a00 + a01 + a11 + a10) / 4)
			dst.Set(x, y, dstColorRGBA64)
		}
	}
}
