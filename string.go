// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
)

const (
	TerminalColor_WHITE = "\033[47m  \033[0m"
	TerminalColor_BLACK = "\033[40m  \033[0m"
)

type StringColor struct {
	Color  string
	ToRGBA func(s string) (r, g, b, a uint32)
}

func (c StringColor) RGBA() (r, g, b, a uint32) {
	if c.ToRGBA != nil {
		return c.ToRGBA(c.Color)
	}
	return color.Gray{0xFF}.RGBA()
}

type StringColorModel struct {
	FromRGBA func(r, g, b, a uint32) string
	ToRGBA   func(s string) (r, g, b, a uint32)
}

func MakeTerminalGrayModel() StringColorModel {
	return StringColorModel{
		FromRGBA: func(r, g, b, a uint32) string {
			if (r + g + b) > 0 {
				return TerminalColor_WHITE
			} else {
				return TerminalColor_BLACK
			}
		},
		ToRGBA: func(s string) (r, g, b, a uint32) {
			switch s {
			case TerminalColor_WHITE:
				return color.Gray{Y: 255}.RGBA()
			case TerminalColor_BLACK:
				return color.Gray{Y: 0}.RGBA()
			default:
				return color.Gray{Y: 255}.RGBA()
			}
		},
	}
}

func MakeJsonColorModel() StringColorModel {
	return StringColorModel{
		FromRGBA: func(r, g, b, a uint32) string {
			var rgba = [...]uint32{r, g, b, a}
			if s, err := json.MarshalIndent(rgba, "", "\t"); err == nil {
				return string(s)
			}
			return ""
		},
		ToRGBA: func(s string) (r, g, b, a uint32) {
			var rgba []uint32
			if err := json.Unmarshal([]byte(s), &rgba); err != nil {
				return
			}
			if len(rgba) > 0 {
				r = rgba[0]
			}
			if len(rgba) > 1 {
				g = rgba[1]
			}
			if len(rgba) > 2 {
				b = rgba[2]
			}
			if len(rgba) > 3 {
				a = rgba[3]
			}
			return
		},
	}
}

func (m StringColorModel) Convert(c color.Color) color.Color {
	return StringColor{
		Color:  m.FromRGBA(c.RGBA()),
		ToRGBA: m.ToRGBA,
	}
}

type StringImage struct {
	Pix    []string
	Rect   image.Rectangle
	Stride int

	StringColorModel
}

func (p *StringImage) Bounds() image.Rectangle { return p.Rect }
func (p *StringImage) Channels() int           { return 1 }
func (p *StringImage) ColorModel() color.Model { return p.StringColorModel }

func (p *StringImage) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return StringColor{}
	}
	i := p.PixOffset(x, y)
	return StringColor{
		Color:  p.Pix[i],
		ToRGBA: p.ToRGBA,
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *StringImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x - p.Rect.Min.X)
}

func (p *StringImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := p.StringColorModel.Convert(c).(StringColor)
	p.Pix[i] = c1.Color
	return
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *StringImage) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &StringImage{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &StringImage{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *StringImage) Opaque() bool {
	return true
}

func (m *StringImage) ToString(zeroColor string) string {
	var b = m.Bounds()
	var buf bytes.Buffer

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := m.At(x, y).(StringColor)
			if c.Color != "" {
				fmt.Print(&buf, c.Color)
			} else {
				fmt.Print(&buf, zeroColor)
			}
		}
		fmt.Print(&buf, "\n")
	}

	return buf.String()
}

// NewStringImage returns a new StringImage with the given bounds.
func NewStringImage(r image.Rectangle, colorModel StringColorModel) *StringImage {
	w, h := r.Dx(), r.Dy()
	pix := make([]string, w*h)

	return &StringImage{
		Pix:    pix,
		Stride: w,
		Rect:   r,

		StringColorModel: colorModel,
	}
}
