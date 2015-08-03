// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"
	"image/draw"

	ximage "github.com/chai2010/image"
)

func MakePyrDown(m image.Image, scaler Scaler) image.Image {
	dst := newPyrDownImage(m)
	dr := dst.Bounds()
	scaler.Scale(dst, dr, m, image.Rect(dr.Min.X, dr.Min.Y, dr.Min.X+dr.Dx()*2, dr.Min.Y+dr.Dy()*2))
	return dst
}

func newPyrDownImage(m image.Image) draw.Image {
	b := m.Bounds()
	r := image.Rect(b.Min.X, b.Min.Y, b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2)
	switch m.(type) {
	case *image.Gray:
		return image.NewGray(r)
	case *image.Gray16:
		return image.NewGray16(r)
	case *image.RGBA:
		return image.NewRGBA(r)
	case *image.RGBA64:
		return image.NewRGBA64(r)
	case *image.YCbCr:
		return image.NewRGBA(r)
	}

	if m, ok := ximage.AsMemPImage(m); ok {
		return ximage.NewMemPImage(r, m.XChannels, m.XDataType)
	}

	// unknown
	return image.NewRGBA64(r)
}
