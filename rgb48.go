// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"image/color"
)

var (
	_ color.Color = (*RGB48Color)(nil)
	_ image.Image = (*RGB48Image)(nil)
)

var RGB48Model color.Model = color.ModelFunc(rgb48Model)

type RGB48Color struct {
	R, G, B uint16
}

func (c RGB48Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	g = uint32(c.G)
	b = uint32(c.B)
	a = 0xffff
	return
}

func rgb48Model(c color.Color) color.Color {
	if _, ok := c.(RGB48Color); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return RGB48Color{R: uint16(r), G: uint16(g), B: uint16(b)}
}

type RGB48Image struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func (p *RGB48Image) ColorModel() color.Model { return color.RGBAModel }

func (p *RGB48Image) Bounds() image.Rectangle { return p.Rect }

func (p *RGB48Image) At(x, y int) color.Color {
	return p.RGB48At(x, y)
}

func (p *RGB48Image) RGB48At(x, y int) RGB48Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return RGB48Color{}
	}
	i := p.PixOffset(x, y)
	return RGB48Color{
		R: uint16(p.Pix[i+0])<<8 | uint16(p.Pix[i+1]),
		G: uint16(p.Pix[i+2])<<8 | uint16(p.Pix[i+3]),
		B: uint16(p.Pix[i+4])<<8 | uint16(p.Pix[i+5]),
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGB48Image) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *RGB48Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := RGB48Model.Convert(c).(RGB48Color)
	p.Pix[i+0] = uint8(c1.R >> 8)
	p.Pix[i+1] = uint8(c1.R)
	p.Pix[i+2] = uint8(c1.G >> 8)
	p.Pix[i+3] = uint8(c1.G)
	p.Pix[i+4] = uint8(c1.B >> 8)
	p.Pix[i+5] = uint8(c1.B)
	return
}

func (p *RGB48Image) SetRGB48(x, y int, c RGB48Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i+0] = uint8(c.R >> 8)
	p.Pix[i+1] = uint8(c.R)
	p.Pix[i+2] = uint8(c.G >> 8)
	p.Pix[i+3] = uint8(c.G)
	p.Pix[i+4] = uint8(c.B >> 8)
	p.Pix[i+5] = uint8(c.B)
	return
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGB48Image) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RGB48Image{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGB48Image{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RGB48Image) Opaque() bool {
	return true
}

// NewRGB48Image returns a new RGB48Image with the given bounds.
func NewRGB48Image(r image.Rectangle) *RGB48Image {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 6*w*h)
	return &RGB48Image{
		Pix:    pix,
		Stride: 6 * w,
		Rect:   r,
	}
}

func NewRGB48ImageFrom(m image.Image) *RGB48Image {
	if m, ok := m.(*RGB48Image); ok {
		return m
	}

	// convert to RGB48Image
	b := m.Bounds()
	rgb := NewRGB48Image(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pr, pg, pb, _ := m.At(x, y).RGBA()
			rgb.SetRGB48(x, y, RGB48Color{
				R: uint16(pr),
				G: uint16(pg),
				B: uint16(pb),
			})
		}
	}
	return rgb
}
