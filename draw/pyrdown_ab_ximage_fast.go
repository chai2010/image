// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"
	"reflect"

	ximage "github.com/chai2010/image"
)

func abPyrDown_xImage(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	if dst.Channels != src.Channels || dst.DataType != src.DataType {
		abPyrDownImage(dst, r, src, sp)
		return
	}

	switch dst.DataType {
	case reflect.Int8:
		abPyrDown_xImage_int8(dst, r, src, sp)
		return
	case reflect.Int16:
		abPyrDown_xImage_int16(dst, r, src, sp)
		return
	case reflect.Int32:
		abPyrDown_xImage_int32(dst, r, src, sp)
		return
	case reflect.Int64:
		abPyrDown_xImage_int64(dst, r, src, sp)
		return
	case reflect.Uint8:
		abPyrDown_xImage_uint8(dst, r, src, sp)
		return
	case reflect.Uint16:
		abPyrDown_xImage_uint16(dst, r, src, sp)
		return
	case reflect.Uint32:
		abPyrDown_xImage_uint32(dst, r, src, sp)
		return
	case reflect.Uint64:
		abPyrDown_xImage_uint64(dst, r, src, sp)
		return
	case reflect.Float32:
		abPyrDown_xImage_float32(dst, r, src, sp)
		return
	case reflect.Float64:
		abPyrDown_xImage_float64(dst, r, src, sp)
		return
	case reflect.Complex64:
		abPyrDown_xImage_complex64(dst, r, src, sp)
		return
	case reflect.Complex128:
		abPyrDown_xImage_complex128(dst, r, src, sp)
		return
	}

	abPyrDownImage(dst, r, src, sp)
	return
}

func abPyrDown_xImage_int8(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*1*dst.Channels*1].Int8s()
		srcLine0 := src.Pix[off1:][:r.Dx()*1*dst.Channels*2].Int8s()
		srcLine1 := src.Pix[off1:][:r.Dx()*1*dst.Channels*2].Int8s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := int16(srcLine0[j+0])
				v01 := int16(srcLine0[j+1])
				v10 := int16(srcLine1[j+0])
				v11 := int16(srcLine1[j+1])
				dstLineX[i] = int8((v00 + v01 + v10 + v11) / 4)
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_int16(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*2*dst.Channels*1].Int16s()
		srcLine0 := src.Pix[off1:][:r.Dx()*2*dst.Channels*2].Int16s()
		srcLine1 := src.Pix[off1:][:r.Dx()*2*dst.Channels*2].Int16s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := int32(srcLine0[j+0])
				v01 := int32(srcLine0[j+1])
				v10 := int32(srcLine1[j+0])
				v11 := int32(srcLine1[j+1])
				dstLineX[i] = int16((v00 + v01 + v10 + v11) / 4)
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_int32(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*4*dst.Channels*1].Int32s()
		srcLine0 := src.Pix[off1:][:r.Dx()*4*dst.Channels*2].Int32s()
		srcLine1 := src.Pix[off1:][:r.Dx()*4*dst.Channels*2].Int32s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := int64(srcLine0[j+0])
				v01 := int64(srcLine0[j+1])
				v10 := int64(srcLine1[j+0])
				v11 := int64(srcLine1[j+1])
				dstLineX[i] = int32((v00 + v01 + v10 + v11) / 4)
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_int64(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*8*dst.Channels*1].Int64s()
		srcLine0 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Int64s()
		srcLine1 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Int64s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := int64(srcLine0[j+0])
				v01 := int64(srcLine0[j+1])
				v10 := int64(srcLine1[j+0])
				v11 := int64(srcLine1[j+1])
				dstLineX[i] = int64((v00 + v01 + v10 + v11) / 4)
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_uint8(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*1*dst.Channels*1].Uint8s()
		srcLine0 := src.Pix[off1:][:r.Dx()*1*dst.Channels*2].Uint8s()
		srcLine1 := src.Pix[off1:][:r.Dx()*1*dst.Channels*2].Uint8s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := uint16(srcLine0[j+0])
				v01 := uint16(srcLine0[j+1])
				v10 := uint16(srcLine1[j+0])
				v11 := uint16(srcLine1[j+1])
				dstLineX[i] = uint8((v00 + v01 + v10 + v11) / 4)
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_uint16(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*2*dst.Channels*1].Uint16s()
		srcLine0 := src.Pix[off1:][:r.Dx()*2*dst.Channels*2].Uint16s()
		srcLine1 := src.Pix[off1:][:r.Dx()*2*dst.Channels*2].Uint16s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := uint32(srcLine0[j+0])
				v01 := uint32(srcLine0[j+1])
				v10 := uint32(srcLine1[j+0])
				v11 := uint32(srcLine1[j+1])
				dstLineX[i] = uint16((v00 + v01 + v10 + v11) / 4)
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_uint32(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*4*dst.Channels*1].Uint32s()
		srcLine0 := src.Pix[off1:][:r.Dx()*4*dst.Channels*2].Uint32s()
		srcLine1 := src.Pix[off1:][:r.Dx()*4*dst.Channels*2].Uint32s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := uint64(srcLine0[j+0])
				v01 := uint64(srcLine0[j+1])
				v10 := uint64(srcLine1[j+0])
				v11 := uint64(srcLine1[j+1])
				dstLineX[i] = uint32((v00 + v01 + v10 + v11) / 4)
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_uint64(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*8*dst.Channels*1].Uint64s()
		srcLine0 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Uint64s()
		srcLine1 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Uint64s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := uint64(srcLine0[j+0])
				v01 := uint64(srcLine0[j+1])
				v10 := uint64(srcLine1[j+0])
				v11 := uint64(srcLine1[j+1])
				dstLineX[i] = uint64((v00 + v01 + v10 + v11) / 4)
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_float32(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	if dst.Channels == 1 {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			dstLineX := dst.Pix[off0:][:r.Dx()*4*dst.Channels*1].Float32s()
			srcLine0 := src.Pix[off1:][:r.Dx()*4*dst.Channels*2].Float32s()
			srcLine1 := src.Pix[off1:][:r.Dx()*4*dst.Channels*2].Float32s()

			for i, j := 0, 0; i < len(dstLineX); i, j = i+1, j+2 {
				v00 := float32(srcLine0[j+0])
				v01 := float32(srcLine0[j+1])
				v10 := float32(srcLine1[j+0])
				v11 := float32(srcLine1[j+1])
				dstLineX[i] = float32((v00 + v01 + v10 + v11) / 4)
			}

			off0 += dst.Stride * 1
			off1 += src.Stride * 2
			off2 += src.Stride * 2
		}
	} else {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			dstLineX := dst.Pix[off0:][:r.Dx()*4*dst.Channels*1].Float32s()
			srcLine0 := src.Pix[off1:][:r.Dx()*4*dst.Channels*2].Float32s()
			srcLine1 := src.Pix[off1:][:r.Dx()*4*dst.Channels*2].Float32s()

			for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
				for k := 0; k < dst.Channels; k++ {
					v00 := float32(srcLine0[j+0])
					v01 := float32(srcLine0[j+1])
					v10 := float32(srcLine1[j+0])
					v11 := float32(srcLine1[j+1])
					dstLineX[i] = float32((v00 + v01 + v10 + v11) / 4)
				}
			}

			off0 += dst.Stride * 1
			off1 += src.Stride * 2
			off2 += src.Stride * 2
		}
	}
}

func abPyrDown_xImage_float64(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	if dst.Channels == 1 {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			dstLineX := dst.Pix[off0:][:r.Dx()*8*dst.Channels*1].Float64s()
			srcLine0 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Float64s()
			srcLine1 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Float64s()

			for i, j := 0, 0; i < len(dstLineX); i, j = i+1, j+2 {
				v00 := float64(srcLine0[j+0])
				v01 := float64(srcLine0[j+1])
				v10 := float64(srcLine1[j+0])
				v11 := float64(srcLine1[j+1])
				dstLineX[i] = float64((v00 + v01 + v10 + v11) / 4)
			}

			off0 += dst.Stride * 1
			off1 += src.Stride * 2
			off2 += src.Stride * 2
		}
	} else {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			dstLineX := dst.Pix[off0:][:r.Dx()*8*dst.Channels*1].Float64s()
			srcLine0 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Float64s()
			srcLine1 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Float64s()

			for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
				for k := 0; k < dst.Channels; k++ {
					v00 := float64(srcLine0[j+0])
					v01 := float64(srcLine0[j+1])
					v10 := float64(srcLine1[j+0])
					v11 := float64(srcLine1[j+1])
					dstLineX[i] = float64((v00 + v01 + v10 + v11) / 4)
				}
			}

			off0 += dst.Stride * 1
			off1 += src.Stride * 2
			off2 += src.Stride * 2
		}
	}
}

func abPyrDown_xImage_complex64(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*8*dst.Channels*1].Complex64s()
		srcLine0 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Complex64s()
		srcLine1 := src.Pix[off1:][:r.Dx()*8*dst.Channels*2].Complex64s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := srcLine0[j+0]
				v01 := srcLine0[j+1]
				v10 := srcLine1[j+0]
				v11 := srcLine1[j+1]
				dstLineX[i] = (v00 + v01 + v10 + v11) / 4
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}

func abPyrDown_xImage_complex128(dst *ximage.Image, r image.Rectangle, src *ximage.Image, sp image.Point) {
	off0 := dst.PixOffset(r.Min.X, r.Min.Y)
	off1 := src.PixOffset(sp.X, sp.Y)
	off2 := off1 + src.Stride

	for y := r.Min.Y; y < r.Max.Y; y++ {
		dstLineX := dst.Pix[off0:][:r.Dx()*16*dst.Channels*1].Complex128s()
		srcLine0 := src.Pix[off1:][:r.Dx()*16*dst.Channels*2].Complex128s()
		srcLine1 := src.Pix[off1:][:r.Dx()*16*dst.Channels*2].Complex128s()

		for i, j := 0, 0; i < len(dstLineX); i, j = i+dst.Channels*1, j+dst.Channels*2 {
			for k := 0; k < dst.Channels; k++ {
				v00 := srcLine0[j+0]
				v01 := srcLine0[j+1]
				v10 := srcLine1[j+0]
				v11 := srcLine1[j+1]
				dstLineX[i] = (v00 + v01 + v10 + v11) / 4
			}
		}

		off0 += dst.Stride * 1
		off1 += src.Stride * 2
		off2 += src.Stride * 2
	}
}
