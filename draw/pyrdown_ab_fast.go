// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"
	"image/color"
	"image/draw"

	ximage "github.com/chai2010/image"
	xdraw "golang.org/x/image/draw"
)

func abPyrDownFast(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	switch dst := dst.(type) {
	case *image.Gray:
		switch src := src.(type) {
		case *image.Gray:
			abPyrDown_Gray_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			abPyrDown_Gray_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			abPyrDown_Gray_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			abPyrDown_Gray_RGBA64(dst, r, src, sp)
			return
		}
	case *image.Gray16:
		switch src := src.(type) {
		case *image.Gray:
			abPyrDown_Gray16_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			abPyrDown_Gray16_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			abPyrDown_Gray16_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			abPyrDown_Gray16_RGBA64(dst, r, src, sp)
			return
		}
	case *ximage.RGBImage:
		switch src := src.(type) {
		case *ximage.RGBImage:
			abPyrDown_RGB_RGB(dst, r, src, sp)
			return
		}
	case *ximage.RGB48Image:
		switch src := src.(type) {
		case *ximage.RGB48Image:
			abPyrDown_RGB48_RGB48(dst, r, src, sp)
			return
		}
	case *image.RGBA:
		switch src := src.(type) {
		case *image.Gray:
			abPyrDown_RGBA_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			abPyrDown_RGBA_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			abPyrDown_RGBA_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			abPyrDown_RGBA_RGBA64(dst, r, src, sp)
			return
		}
	case *image.RGBA64:
		switch src := src.(type) {
		case *image.Gray:
			abPyrDown_RGBA64_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			abPyrDown_RGBA64_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			abPyrDown_RGBA64_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			abPyrDown_RGBA64_RGBA64(dst, r, src, sp)
			return
		}
	}

	if dst, ok := ximage.AsMemPImage(dst); ok {
		if src, ok := ximage.AsMemPImage(src); ok {
			abPyrDown_xImage(dst, r, src, sp)
			return
		}
	}

	abPyrDownImage(dst, r, src, sp)
}

func abPyrDown_Gray_Gray(dst *image.Gray, r image.Rectangle, src *image.Gray, sp image.Point) {
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

func abPyrDown_Gray_Gray16(dst *image.Gray, r image.Rectangle, src *image.Gray16, sp image.Point) {
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

func abPyrDown_Gray_RGBA(dst *image.Gray, r image.Rectangle, src *image.RGBA, sp image.Point) {
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

func abPyrDown_Gray_RGBA64(dst *image.Gray, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*1*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*8*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*8*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+1*1, j+8*2 {
			r00 := uint16(srcLine0[j+8*0+0])<<8 | uint16(srcLine0[j+8*0+1])
			g00 := uint16(srcLine0[j+8*0+2])<<8 | uint16(srcLine0[j+8*0+3])
			b00 := uint16(srcLine0[j+8*0+4])<<8 | uint16(srcLine0[j+8*0+5])

			r01 := uint16(srcLine0[j+8*1+0])<<8 | uint16(srcLine0[j+8*1+1])
			g01 := uint16(srcLine0[j+8*1+2])<<8 | uint16(srcLine0[j+8*1+3])
			b01 := uint16(srcLine0[j+8*1+4])<<8 | uint16(srcLine0[j+8*1+5])

			r10 := uint16(srcLine1[j+8*0+0])<<8 | uint16(srcLine1[j+8*0+1])
			g10 := uint16(srcLine1[j+8*0+2])<<8 | uint16(srcLine1[j+8*0+3])
			b10 := uint16(srcLine1[j+8*0+4])<<8 | uint16(srcLine1[j+8*0+5])

			r11 := uint16(srcLine1[j+8*1+0])<<8 | uint16(srcLine1[j+8*1+1])
			g11 := uint16(srcLine1[j+8*1+2])<<8 | uint16(srcLine1[j+8*1+3])
			b11 := uint16(srcLine1[j+8*1+4])<<8 | uint16(srcLine1[j+8*1+5])

			rxx := (uint32(r00) + uint32(r01) + uint32(r10) + uint32(r11)) >> 2
			gxx := (uint32(g00) + uint32(g01) + uint32(g10) + uint32(g11)) >> 2
			bxx := (uint32(b00) + uint32(b01) + uint32(b10) + uint32(b11)) >> 2

			dstLineX[i] = uint8(((299*rxx + 587*gxx + 114*bxx + 500) / 1000) >> 8)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_Gray16_Gray(dst *image.Gray16, r image.Rectangle, src *image.Gray, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*2*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*1*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*1*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+2*1, j+1*2 {
			v00 := srcLine0[j+0]
			v01 := srcLine0[j+1]
			v10 := srcLine1[j+0]
			v11 := srcLine1[j+1]

			vxx := ((uint32(v00) + uint32(v01) + uint32(v10) + uint32(v11)) * 0x101) >> 2
			dstLineX[i+0] = uint8(vxx >> 8)
			dstLineX[i+1] = uint8(vxx)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_Gray16_Gray16(dst *image.Gray16, r image.Rectangle, src *image.Gray16, sp image.Point) {
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

			vxx := uint16((v00 + v01 + v10 + v11) >> 2)
			dstLineX[i+0] = uint8(vxx >> 8)
			dstLineX[i+1] = uint8(vxx)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_Gray16_RGBA(dst *image.Gray16, r image.Rectangle, src *image.RGBA, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*2*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*4*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*4*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+2*1, j+4*2 {
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

			vxx := uint16((299*rxx + 587*gxx + 114*bxx + 500) / 1000)
			dstLineX[i+0] = uint8(vxx >> 8)
			dstLineX[i+1] = uint8(vxx)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_Gray16_RGBA64(dst *image.Gray16, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*2*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*8*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*8*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+2*1, j+8*2 {
			r00 := uint16(srcLine0[j+8*0+0])<<8 | uint16(srcLine0[j+8*0+1])
			g00 := uint16(srcLine0[j+8*0+2])<<8 | uint16(srcLine0[j+8*0+3])
			b00 := uint16(srcLine0[j+8*0+4])<<8 | uint16(srcLine0[j+8*0+5])

			r01 := uint16(srcLine0[j+8*1+0])<<8 | uint16(srcLine0[j+8*1+1])
			g01 := uint16(srcLine0[j+8*1+2])<<8 | uint16(srcLine0[j+8*1+3])
			b01 := uint16(srcLine0[j+8*1+4])<<8 | uint16(srcLine0[j+8*1+5])

			r10 := uint16(srcLine1[j+8*0+0])<<8 | uint16(srcLine1[j+8*0+1])
			g10 := uint16(srcLine1[j+8*0+2])<<8 | uint16(srcLine1[j+8*0+3])
			b10 := uint16(srcLine1[j+8*0+4])<<8 | uint16(srcLine1[j+8*0+5])

			r11 := uint16(srcLine1[j+8*1+0])<<8 | uint16(srcLine1[j+8*1+1])
			g11 := uint16(srcLine1[j+8*1+2])<<8 | uint16(srcLine1[j+8*1+3])
			b11 := uint16(srcLine1[j+8*1+4])<<8 | uint16(srcLine1[j+8*1+5])

			rxx := (uint32(r00) + uint32(r01) + uint32(r10) + uint32(r11)) >> 2
			gxx := (uint32(g00) + uint32(g01) + uint32(g10) + uint32(g11)) >> 2
			bxx := (uint32(b00) + uint32(b01) + uint32(b10) + uint32(b11)) >> 2

			vxx := uint16((299*rxx + 587*gxx + 114*bxx + 500) / 1000)
			dstLineX[i+0] = uint8(vxx >> 8)
			dstLineX[i+1] = uint8(vxx)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_RGB_RGB(dst *ximage.RGBImage, r image.Rectangle, src *ximage.RGBImage, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.XStride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.XPix[off0:][:r.Dx()*3*1]
		srcLine0 := src.XPix[off1:][:r.Dx()*3*2]
		srcLine1 := src.XPix[off2:][:r.Dx()*3*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+3*1, j+3*2 {
			r00 := srcLine0[j+0]
			g00 := srcLine0[j+1]
			b00 := srcLine0[j+2]

			r01 := srcLine0[j+3]
			g01 := srcLine0[j+4]
			b01 := srcLine0[j+5]

			r10 := srcLine1[j+0]
			g10 := srcLine1[j+1]
			b10 := srcLine1[j+2]

			r11 := srcLine1[j+3]
			g11 := srcLine1[j+4]
			b11 := srcLine1[j+5]

			dstLineX[i+0] = uint8(((uint32(r00) + uint32(r01) + uint32(r10) + uint32(r11)) * 0x101) >> 10)
			dstLineX[i+1] = uint8(((uint32(g00) + uint32(g01) + uint32(g10) + uint32(g11)) * 0x101) >> 10)
			dstLineX[i+2] = uint8(((uint32(b00) + uint32(b01) + uint32(b10) + uint32(b11)) * 0x101) >> 10)
		}

		off0 += dst.XStride * 1
		off1 += src.XStride * 2
		off2 += src.XStride * 2
	}
}

func abPyrDown_RGB48_RGB48(dst *ximage.RGB48Image, r image.Rectangle, src *ximage.RGB48Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.XStride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.XPix[off0:][:r.Dx()*6*1]
		srcLine0 := src.XPix[off1:][:r.Dx()*6*2]
		srcLine1 := src.XPix[off2:][:r.Dx()*6*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+6*1, j+6*2 {
			r00 := uint32(srcLine0[j+8*0+0])<<8 | uint32(srcLine0[j+8*0+1])
			g00 := uint32(srcLine0[j+8*0+2])<<8 | uint32(srcLine0[j+8*0+3])
			b00 := uint32(srcLine0[j+8*0+4])<<8 | uint32(srcLine0[j+8*0+5])
			r01 := uint32(srcLine0[j+8*1+0])<<8 | uint32(srcLine0[j+8*1+1])
			g01 := uint32(srcLine0[j+8*1+2])<<8 | uint32(srcLine0[j+8*1+3])
			b01 := uint32(srcLine0[j+8*1+4])<<8 | uint32(srcLine0[j+8*1+5])
			r10 := uint32(srcLine1[j+8*0+0])<<8 | uint32(srcLine1[j+8*0+1])
			g10 := uint32(srcLine1[j+8*0+2])<<8 | uint32(srcLine1[j+8*0+3])
			b10 := uint32(srcLine1[j+8*0+4])<<8 | uint32(srcLine1[j+8*0+5])
			r11 := uint32(srcLine1[j+8*1+0])<<8 | uint32(srcLine1[j+8*1+1])
			g11 := uint32(srcLine1[j+8*1+2])<<8 | uint32(srcLine1[j+8*1+3])
			b11 := uint32(srcLine1[j+8*1+4])<<8 | uint32(srcLine1[j+8*1+5])

			rxx := uint16((r00 + r01 + r10 + r11) >> 2)
			gxx := uint16((g00 + g01 + g10 + g11) >> 2)
			bxx := uint16((b00 + b01 + b10 + b11) >> 2)

			dstLineX[i+0] = uint8(rxx >> 8)
			dstLineX[i+1] = uint8(rxx)
			dstLineX[i+2] = uint8(gxx >> 8)
			dstLineX[i+3] = uint8(gxx)
			dstLineX[i+4] = uint8(bxx >> 8)
			dstLineX[i+5] = uint8(bxx)
		}

		off0 += dst.XStride * 1
		off1 += src.XStride * 2
		off2 += src.XStride * 2
	}
}

func abPyrDown_RGBA_Gray(dst *image.RGBA, r image.Rectangle, src *image.Gray, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*4*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*1*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*1*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+4*1, j+1*2 {
			v00 := srcLine0[j+0]
			v01 := srcLine0[j+1]
			v10 := srcLine1[j+0]
			v11 := srcLine1[j+1]

			vxx := uint8(((uint32(v00) + uint32(v01) + uint32(v10) + uint32(v11)) * 0x101) >> 10)
			dstLineX[i+0] = vxx
			dstLineX[i+1] = vxx
			dstLineX[i+2] = vxx
			dstLineX[i+3] = 0xff
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_RGBA_Gray16(dst *image.RGBA, r image.Rectangle, src *image.Gray16, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*4*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*2*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*2*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+4*1, j+2*2 {
			v00 := uint16(srcLine0[j+0])<<8 | uint16(srcLine0[j+1])
			v01 := uint16(srcLine0[j+2])<<8 | uint16(srcLine0[j+3])
			v10 := uint16(srcLine1[j+0])<<8 | uint16(srcLine0[j+1])
			v11 := uint16(srcLine1[j+2])<<8 | uint16(srcLine0[j+3])

			vxx := uint8((uint32(v00) + uint32(v01) + uint32(v10) + uint32(v11)) >> 10)
			dstLineX[i+0] = vxx
			dstLineX[i+1] = vxx
			dstLineX[i+2] = vxx
			dstLineX[i+3] = 0xff
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_RGBA_RGBA(dst *image.RGBA, r image.Rectangle, src *image.RGBA, sp image.Point) {
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

func abPyrDown_RGBA_RGBA64(dst *image.RGBA, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*4*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*8*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*8*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+4*1, j+8*2 {
			r00 := uint16(srcLine0[j+8*0+0])<<8 | uint16(srcLine0[j+8*0+1])
			g00 := uint16(srcLine0[j+8*0+2])<<8 | uint16(srcLine0[j+8*0+3])
			b00 := uint16(srcLine0[j+8*0+4])<<8 | uint16(srcLine0[j+8*0+5])
			a00 := uint16(srcLine0[j+8*0+6])<<8 | uint16(srcLine0[j+8*0+7])

			r01 := uint16(srcLine0[j+8*1+0])<<8 | uint16(srcLine0[j+8*1+1])
			g01 := uint16(srcLine0[j+8*1+2])<<8 | uint16(srcLine0[j+8*1+3])
			b01 := uint16(srcLine0[j+8*1+4])<<8 | uint16(srcLine0[j+8*1+5])
			a01 := uint16(srcLine0[j+8*1+6])<<8 | uint16(srcLine0[j+8*1+7])

			r10 := uint16(srcLine1[j+8*0+0])<<8 | uint16(srcLine1[j+8*0+1])
			g10 := uint16(srcLine1[j+8*0+2])<<8 | uint16(srcLine1[j+8*0+3])
			b10 := uint16(srcLine1[j+8*0+4])<<8 | uint16(srcLine1[j+8*0+5])
			a10 := uint16(srcLine1[j+8*0+6])<<8 | uint16(srcLine1[j+8*0+7])

			r11 := uint16(srcLine1[j+8*1+0])<<8 | uint16(srcLine1[j+8*1+1])
			g11 := uint16(srcLine1[j+8*1+2])<<8 | uint16(srcLine1[j+8*1+3])
			b11 := uint16(srcLine1[j+8*1+4])<<8 | uint16(srcLine1[j+8*1+5])
			a11 := uint16(srcLine1[j+8*1+6])<<8 | uint16(srcLine1[j+8*1+7])

			rxx := (uint32(r00) + uint32(r01) + uint32(r10) + uint32(r11)) >> 2
			gxx := (uint32(g00) + uint32(g01) + uint32(g10) + uint32(g11)) >> 2
			bxx := (uint32(b00) + uint32(b01) + uint32(b10) + uint32(b11)) >> 2
			axx := (uint32(a00) + uint32(a01) + uint32(a10) + uint32(a11)) >> 2

			dstLineX[i+0] = uint8(rxx >> 8)
			dstLineX[i+1] = uint8(gxx >> 8)
			dstLineX[i+2] = uint8(bxx >> 8)
			dstLineX[i+3] = uint8(axx >> 8)
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_RGBA64_Gray(dst *image.RGBA64, r image.Rectangle, src *image.Gray, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*4*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*1*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*1*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+4*1, j+1*2 {
			v00 := srcLine0[j+0]
			v01 := srcLine0[j+1]
			v10 := srcLine1[j+0]
			v11 := srcLine1[j+1]

			vxx := uint16(((uint32(v00) + uint32(v01) + uint32(v10) + uint32(v11)) * 0x101) >> 2)
			dstLineX[i+0] = uint8(vxx >> 8)
			dstLineX[i+1] = uint8(vxx)
			dstLineX[i+2] = uint8(vxx >> 8)
			dstLineX[i+3] = uint8(vxx)
			dstLineX[i+4] = uint8(vxx >> 8)
			dstLineX[i+5] = uint8(vxx)
			dstLineX[i+6] = 0xff
			dstLineX[i+7] = 0xff
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_RGBA64_Gray16(dst *image.RGBA64, r image.Rectangle, src *image.Gray16, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*4*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*2*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*2*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+4*1, j+2*2 {
			v00 := uint32(srcLine0[j+0])<<8 | uint32(srcLine0[j+1])
			v01 := uint32(srcLine0[j+2])<<8 | uint32(srcLine0[j+3])
			v10 := uint32(srcLine1[j+0])<<8 | uint32(srcLine1[j+1])
			v11 := uint32(srcLine1[j+2])<<8 | uint32(srcLine1[j+3])

			vxx := uint16((uint32(v00) + uint32(v01) + uint32(v10) + uint32(v11)) >> 2)
			dstLineX[i+0] = uint8(vxx >> 8)
			dstLineX[i+1] = uint8(vxx)
			dstLineX[i+2] = uint8(vxx >> 8)
			dstLineX[i+3] = uint8(vxx)
			dstLineX[i+4] = uint8(vxx >> 8)
			dstLineX[i+5] = uint8(vxx)
			dstLineX[i+6] = 0xff
			dstLineX[i+7] = 0xff
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_RGBA64_RGBA(dst *image.RGBA64, r image.Rectangle, src *image.RGBA, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*8*1]
		srcLine0 := src.Pix[off1:][:r.Dx()*4*2]
		srcLine1 := src.Pix[off2:][:r.Dx()*4*2]

		for i, j := 0, 0; i < len(dstLineX); i, j = i+8*1, j+4*2 {
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

			rxx := uint16(((uint32(r00) + uint32(r01) + uint32(r10) + uint32(r11)) * 0x101) >> 2)
			gxx := uint16(((uint32(g00) + uint32(g01) + uint32(g10) + uint32(g11)) * 0x101) >> 2)
			bxx := uint16(((uint32(b00) + uint32(b01) + uint32(b10) + uint32(b11)) * 0x101) >> 2)
			axx := uint16(((uint32(a00) + uint32(a01) + uint32(a10) + uint32(a11)) * 0x101) >> 2)

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

	xdraw.ApproxBiLinear.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		xdraw.Src, nil,
	)
}

func abPyrDown_RGBA64_RGBA64(dst *image.RGBA64, r image.Rectangle, src *image.RGBA64, sp image.Point) {
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

			rxx := uint16((r00 + r01 + r10 + r11) >> 2)
			gxx := uint16((g00 + g01 + g10 + g11) >> 2)
			bxx := uint16((b00 + b01 + b10 + b11) >> 2)
			axx := uint16((a00 + a01 + a10 + a11) >> 2)

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

func abPyrDownImage(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	dstColorRGBA64 := &color.RGBA64{}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			sx := sp.X + (x-r.Min.X)*2
			sy := sp.Y + (y-r.Min.Y)*2

			r00, g00, b00, a00 := src.At(sx+0, sy+0).RGBA()
			r01, g01, b01, a01 := src.At(sx+0, sy+1).RGBA()
			r11, g11, b11, a11 := src.At(sx+1, sy+1).RGBA()
			r10, g10, b10, a10 := src.At(sx+1, sy+0).RGBA()

			dstColorRGBA64.R = uint16((r00 + r01 + r11 + r10) >> 2)
			dstColorRGBA64.G = uint16((g00 + g01 + g11 + g10) >> 2)
			dstColorRGBA64.B = uint16((b00 + b01 + b11 + b10) >> 2)
			dstColorRGBA64.A = uint16((a00 + a01 + a11 + a10) >> 2)
			dst.Set(x, y, dstColorRGBA64)
		}
	}
}
