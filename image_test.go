// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"bytes"
	"image"
	"reflect"
	"testing"
	"unsafe"
)

type tUnknown struct {
	MemPMagic string // MemP
	Rect      image.Rectangle
	_Channels int
	_DataType reflect.Kind
	Pix       []float32

	// Stride is the Pix stride (in bytes, must align with PixelSize)
	// between vertically adjacent pixels.
	_Stride int
}

func (p *tUnknown) Equal(q *Image) bool {
	if p.MemPMagic != q.MemPMagic {
		return false
	}
	if p.Rect != q.Rect {
		return false
	}
	if p._Channels != q.Channels {
		return false
	}
	if p._DataType != q.DataType {
		return false
	}
	if !bytes.Equal(AsPixSilce(p.Pix), q.Pix) {
		return false
	}
	if p._Stride != q.Stride {
		return false
	}
	return true
}

func TestImage(t *testing.T) {
	//
}

func TestImage_Unsafe(t *testing.T) {
	b := image.Rect(0, 0, 300, 400)
	m1 := (*tUnknown)(unsafe.Pointer(NewImage(b, 3, reflect.Uint8)))
	m2, _ := AsMemPImage(m1)
	m3, _ := AsMemPImage(*m1)

	if !m1.Equal(m2) {
		m1.Pix, m2.Pix = nil, nil
		t.Fatalf("not equal: %v != %v", m1, m2)
	}
	if !m1.Equal(m3) {
		m1.Pix, m3.Pix = nil, nil
		t.Fatalf("not equal: %v != %v", m1, m3)
	}
}
