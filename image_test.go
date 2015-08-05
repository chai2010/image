// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"bytes"
	"image"
	"reflect"
	"testing"
)

var (
	_ MemP = (*tUnknown)(nil)
)

type tUnknownPix []byte

type tUnknown struct {
	XMemPMagic string // MemP
	XRect      image.Rectangle
	XChannels  int
	XDataType  reflect.Kind
	XPix       tUnknownPix

	// Stride is the Pix stride (in bytes, must align with PixelSize)
	// between vertically adjacent pixels.
	XStride int
}

func tNewUnknown(r image.Rectangle, channels int, dataType reflect.Kind) *tUnknown {
	m := &tUnknown{
		XMemPMagic: MemPMagic,
		XRect:      r,
		XStride:    r.Dx() * channels * SizeofKind(dataType),
		XChannels:  channels,
		XDataType:  dataType,
	}
	m.XPix = make([]byte, r.Dy()*m.XStride)
	return m
}

func (p *tUnknown) MemPMagic() string {
	return p.XMemPMagic
}

func (p *tUnknown) Bounds() image.Rectangle {
	return p.XRect
}

func (p *tUnknown) Channels() int {
	return p.XChannels
}

func (p *tUnknown) DataType() reflect.Kind {
	return p.XDataType
}

func (p *tUnknown) Pix() []byte {
	return p.XPix
}

func (p *tUnknown) Stride() int {
	return p.XStride
}

func (p *tUnknown) Equal(q *MemPImage) bool {
	if p.XMemPMagic != q.XMemPMagic {
		return false
	}
	if p.XRect != q.XRect {
		return false
	}
	if p.XChannels != q.XChannels {
		return false
	}
	if p.XDataType != q.XDataType {
		return false
	}
	if !bytes.Equal(p.XPix, q.XPix) {
		return false
	}
	if p.XStride != q.XStride {
		return false
	}
	return true
}

func TestImage(t *testing.T) {
	//
}

func TestImage_Unsafe(t *testing.T) {
	b := image.Rect(0, 0, 300, 400)
	m1 := tNewUnknown(b, 3, reflect.Uint8)
	m2 := NewMemPImage(b, 3, reflect.Uint8)

	if !m1.Equal(m2) {
		m1.XPix, m2.XPix = nil, nil
		t.Fatalf("not equal: %v != %v", m1, m2)
	}
}
