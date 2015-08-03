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

func drawFast(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	switch dst := dst.(type) {
	case *image.Gray:
		switch src := src.(type) {
		case *image.Gray:
			drawGray_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			drawGray_Gray16(dst, r, src, sp)
			return
		case *image.RGBA:
			drawGray_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			drawGray_RGBA64(dst, r, src, sp)
			return
		case *image.YCbCr:
			drawGray_YCbCr(dst, r, src, sp)
			return
		}
	case *image.Gray16:
		switch src := src.(type) {
		case *image.Gray16:
			drawGray16_Gray16(dst, r, src, sp)
			return
		case *image.Gray:
			drawGray16_Gray(dst, r, src, sp)
			return
		case *image.YCbCr:
			drawGray16_YCbCr(dst, r, src, sp)
			return
		}
	case *ximage.RGBImage:
		switch src := src.(type) {
		case *ximage.RGBImage:
			drawRGB_RGB(dst, r, src, sp)
			return
		}
	case *ximage.RGB48Image:
		switch src := src.(type) {
		case *ximage.RGB48Image:
			drawRGB48_RGB48(dst, r, src, sp)
			return
		}
	case *image.RGBA:
		switch src := src.(type) {
		case *image.RGBA:
			drawRGBA_RGBA(dst, r, src, sp)
			return
		case *image.RGBA64:
			drawRGBA_RGBA64(dst, r, src, sp)
			return
		case *image.Gray:
			drawRGBA_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			drawRGBA_Gray16(dst, r, src, sp)
			return
		case *image.YCbCr:
			drawRGBA_YCbCr(dst, r, src, sp)
			return
		}
	case *image.RGBA64:
		switch src := src.(type) {
		case *image.RGBA64:
			drawRGBA64_RGBA64(dst, r, src, sp)
			return
		case *image.RGBA:
			drawRGBA64_RGBA(dst, r, src, sp)
			return
		case *image.Gray:
			drawRGBA64_Gray(dst, r, src, sp)
			return
		case *image.Gray16:
			drawRGBA64_Gray16(dst, r, src, sp)
			return
		case *image.YCbCr:
			drawRGBA64_YCbCr(dst, r, src, sp)
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

func drawGray_Gray(dst *image.Gray, r image.Rectangle, src *image.Gray, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*1*1]
		srcLine := src.Pix[off1:][:r.Dx()*1*1]

		copy(dstLine, srcLine)
	}
}

func drawGray_Gray16(dst *image.Gray, r image.Rectangle, src *image.Gray16, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0] = src.Pix[off1]
			off0 += 1
			off1 += 2
		}
	}
}

func drawGray_RGBA(dst *image.Gray, r image.Rectangle, src *image.RGBA, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			rxx := uint32(src.Pix[off1+0]) * 0x101
			gxx := uint32(src.Pix[off1+1]) * 0x101
			bxx := uint32(src.Pix[off1+2]) * 0x101

			dst.Pix[off0] = uint8(((299*rxx + 587*gxx + 114*bxx + 500) / 1000) >> 8)
			off0 += 1
			off1 += 4
		}
	}
}

func drawGray_RGBA64(dst *image.Gray, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			rxx, gxx, bxx, _ := src.RGBA64At(x, y).RGBA()
			dst.Pix[off0] = uint8(((299*rxx + 587*gxx + 114*bxx + 500) / 1000) >> 8)
			off0 += 1
			off1 += 4
		}
	}
}

func drawGray_YCbCr(dst *image.Gray, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.YOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0] = src.Y[off1]
			off0 += 1
			off1 += 1
		}
	}
}

func drawGray16_Gray16(dst *image.Gray16, r image.Rectangle, src *image.Gray16, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*2*1]
		srcLine := src.Pix[off1:][:r.Dx()*2*1]

		copy(dstLine, srcLine)
	}
}

func drawGray16_Gray(dst *image.Gray16, r image.Rectangle, src *image.Gray, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Pix[off1]
			dst.Pix[off0+1] = src.Pix[off1]
			off0 += 2
			off1 += 1
		}
	}
}

func drawGray16_YCbCr(dst *image.Gray16, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.YOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Y[off1]
			dst.Pix[off0+1] = src.Y[off1]
			off0 += 2
			off1 += 1
		}
	}
}

func drawRGB_RGB(dst *ximage.RGBImage, r image.Rectangle, src *ximage.RGBImage, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*1*3]
		srcLine := src.Pix[off1:][:r.Dx()*1*3]

		copy(dstLine, srcLine)
	}
}

func drawRGB48_RGB48(dst *ximage.RGB48Image, r image.Rectangle, src *ximage.RGB48Image, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*2*3]
		srcLine := src.Pix[off1:][:r.Dx()*2*3]

		copy(dstLine, srcLine)
	}
}

func drawRGBA_RGBA(dst *image.RGBA, r image.Rectangle, src *image.RGBA, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*1*4]
		srcLine := src.Pix[off1:][:r.Dx()*1*4]

		copy(dstLine, srcLine)
	}
}

func drawRGBA_RGBA64(dst *image.RGBA, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Pix[off1+0*2]
			dst.Pix[off0+1] = src.Pix[off1+1*2]
			dst.Pix[off0+2] = src.Pix[off1+2*2]
			dst.Pix[off0+3] = src.Pix[off1+3*2]
			off0 += 4
			off1 += 8
		}
	}
}

func drawRGBA_Gray(dst *image.RGBA, r image.Rectangle, src *image.Gray, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Pix[off1]
			dst.Pix[off0+1] = src.Pix[off1]
			dst.Pix[off0+2] = src.Pix[off1]
			dst.Pix[off0+3] = 0xff
			off0 += 4
			off1 += 1
		}
	}
}

func drawRGBA_Gray16(dst *image.RGBA, r image.Rectangle, src *image.Gray16, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Pix[off1]
			dst.Pix[off0+1] = src.Pix[off1]
			dst.Pix[off0+2] = src.Pix[off1]
			dst.Pix[off0+3] = 0xff
			off0 += 4
			off1 += 2
		}
	}
}

func drawRGBA_YCbCr(dst *image.RGBA, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			c := src.YCbCrAt(x, y)
			rxx, gxx, bxx := color.YCbCrToRGB(c.Y, c.Cb, c.Cr)
			dst.Pix[off0+0] = rxx
			dst.Pix[off0+1] = gxx
			dst.Pix[off0+2] = bxx
			dst.Pix[off0+3] = 0xff
			off0 += 4
		}
	}
}

func drawRGBA64_RGBA64(dst *image.RGBA64, r image.Rectangle, src *image.RGBA64, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.Pix[off0:][:r.Dx()*2*4]
		srcLine := src.Pix[off1:][:r.Dx()*2*4]

		copy(dstLine, srcLine)
	}
}

func drawRGBA64_RGBA(dst *image.RGBA64, r image.Rectangle, src *image.RGBA, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Pix[off1+0]
			dst.Pix[off0+1] = src.Pix[off1+0]
			dst.Pix[off0+2] = src.Pix[off1+1]
			dst.Pix[off0+3] = src.Pix[off1+1]
			dst.Pix[off0+4] = src.Pix[off1+2]
			dst.Pix[off0+5] = src.Pix[off1+2]
			dst.Pix[off0+6] = src.Pix[off1+3]
			dst.Pix[off0+7] = src.Pix[off1+3]
			off0 += 8
			off1 += 4
		}
	}
}

func drawRGBA64_Gray(dst *image.RGBA64, r image.Rectangle, src *image.Gray, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Pix[off1]
			dst.Pix[off0+1] = src.Pix[off1]
			dst.Pix[off0+2] = src.Pix[off1]
			dst.Pix[off0+3] = src.Pix[off1]
			dst.Pix[off0+4] = src.Pix[off1]
			dst.Pix[off0+5] = src.Pix[off1]
			dst.Pix[off0+6] = 0xff
			dst.Pix[off0+7] = 0xff
			off0 += 8
			off1 += 1
		}
	}
}

func drawRGBA64_Gray16(dst *image.RGBA64, r image.Rectangle, src *image.Gray16, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Pix[off1]
			dst.Pix[off0+1] = src.Pix[off1]
			dst.Pix[off0+2] = src.Pix[off1]
			dst.Pix[off0+3] = src.Pix[off1]
			dst.Pix[off0+4] = src.Pix[off1]
			dst.Pix[off0+5] = src.Pix[off1]
			dst.Pix[off0+6] = 0xff
			dst.Pix[off0+7] = 0xff
			off0 += 8
			off1 += 2
		}
	}
}

func drawRGBA64_YCbCr(dst *image.RGBA64, r image.Rectangle, src *image.YCbCr, sp image.Point) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.YOffset(sp.X, y)

		for x := r.Min.X; x < r.Max.X; x++ {
			dst.Pix[off0+0] = src.Y[off1]
			dst.Pix[off0+1] = src.Y[off1]
			dst.Pix[off0+2] = src.Y[off1]
			dst.Pix[off0+3] = src.Y[off1]
			dst.Pix[off0+4] = src.Y[off1]
			dst.Pix[off0+5] = src.Y[off1]
			dst.Pix[off0+6] = 0xff
			dst.Pix[off0+7] = 0xff
			off0 += 8
			off1 += 1
		}
	}
}

func drawImage(dst *ximage.MemPImage, r image.Rectangle, src *ximage.MemPImage, sp image.Point) {
	if dst.XChannels != src.XChannels || dst.XDataType != src.XDataType {
		xdraw.Draw(dst, r, src, sp, xdraw.Src)
		return
	}

	dxSize := ximage.SizeofPixel(dst.XChannels, dst.XDataType) * r.Dx()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		off0 := dst.PixOffset(r.Min.X, y)
		off1 := src.PixOffset(sp.X, y)

		dstLine := dst.XPix[off0:][:dxSize]
		srcLine := src.XPix[off1:][:dxSize]

		copy(dstLine, srcLine)
	}
}
