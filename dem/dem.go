// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dem

import (
	"image"
	"image/color"
	"reflect"
)

const (
	InvalidDemValue float32 = -99999 // if d.GetZ(x,y) <= InvalidValue { ... }
)

var (
	_ image.Image = (*Dem)(nil)
)

type Dem Image // Dem is NewImage(r, 1, reflect.Float32 type)

func NewDem(r image.Rectangle) *Dem {
	m := NewImage(r, 1, reflect.Float32)
	return (*Dem)(m)
}

func AsDem(m interface{}) (p *Dem, ok bool) {
	if x, ok := AsMemPImage(m); ok {
		if x.Channels == 1 && x.DataType == reflect.Float32 {
			return (*Dem)(x), true
		}
	}
	return nil, false
}

func (p *Dem) Clone() *Dem {
	q := new(Dem)
	*q = *p
	q.Pix = append([]byte(nil), p.Pix...)
	return q
}

func (p *Dem) Bounds() image.Rectangle {
	return p.Rect
}

func (p *Dem) ColorModel() color.Model {
	return ColorModel(p.Channels, p.DataType)
}

func (p *Dem) At(x, y int) color.Color {
	return (*Image)(p).At(x, y)
}

func (p *Dem) ValueAt(x, y int) float32 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return 0
	}
	i := p.PixOffset(x, y)
	n := SizeofPixel(p.Channels, p.DataType)
	return (p.Pix[i:][:n]).Float32s()[0]
}

func (p *Dem) Set(x, y int, c color.Color) {
	(*Image)(p).Set(x, y, c)
}

func (p *Dem) SetValue(x, y int, c float32) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	n := SizeofPixel(p.Channels, p.DataType)
	(p.Pix[i:][:n]).Float32s()[i] = c
}

func (p *Dem) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*SizeofPixel(p.Channels, p.DataType)
}

func (p *Dem) SubDem(r image.Rectangle) *Dem {
	return (*Dem)((*Image)(p).SubImage(r).(*Image))
}

func (p *Dem) StdImage() image.Image {
	return (*Image)(p).StdImage()
}
