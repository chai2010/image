// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memp

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"testing"

	"golang.org/x/image/draw"
)

func TestPyrDown_gray(t *testing.T) {
	src := tToGray(tLoadImage("./testdata/lena.png"))

	dst0 := image.NewGray(image.Rect(0, 0, src.Bounds().Dx()/2, src.Bounds().Dy()/2))
	dst1 := image.NewGray(image.Rect(0, 0, src.Bounds().Dx()/2, src.Bounds().Dy()/2))

	draw.ApproxBiLinear.Scale(
		dst0, dst0.Bounds(),
		src, src.Bounds(),
		draw.Src, nil,
	)
	_PyrDown_ApproxBiLinear_Gray_Gray(
		dst1, dst1.Bounds(),
		src, image.Pt(0, 0),
	)

	for y := 0; y < src.Bounds().Dy()/2; y++ {
		for x := 0; x < src.Bounds().Dx()/2; x++ {
			if v1, v2 := dst0.GrayAt(x, y), dst1.GrayAt(x, y); v1 != v2 {
				t.Logf("(%d,%d): %v, %v", x*x, 2*y,
					src.GrayAt(x*2+0, y*2+0),
					src.GrayAt(x*2+0, y*2+1),
				)
				t.Logf("(%d,%d): %v, %v", x*x, 2*y,
					src.GrayAt(x*2+1, y*2+1),
					src.GrayAt(x*2+1, y*2+0),
				)
				t.Fatalf("(%d,%d): %v != %v", x, y, v1, v2)
			}
		}
	}
}

func TestPyrDown_gray16(t *testing.T) {
	src := tToGray16(tLoadImage("./testdata/lena.png"))

	dst0 := image.NewGray16(image.Rect(0, 0, src.Bounds().Dx()/2, src.Bounds().Dy()/2))
	dst1 := image.NewGray16(image.Rect(0, 0, src.Bounds().Dx()/2, src.Bounds().Dy()/2))

	draw.ApproxBiLinear.Scale(
		dst0, dst0.Bounds(),
		src, src.Bounds(),
		draw.Src, nil,
	)
	_PyrDown_ApproxBiLinear_Gray16_Gray16(
		dst1, dst1.Bounds(),
		src, image.Pt(0, 0),
	)

	for y := 0; y < src.Bounds().Dy()/2; y++ {
		for x := 0; x < src.Bounds().Dx()/2; x++ {
			if v1, v2 := dst0.Gray16At(x, y), dst1.Gray16At(x, y); v1 != v2 {
				t.Logf("(%d,%d): %x, %x", x*x, 2*y,
					src.Gray16At(x*2+0, y*2+0),
					src.Gray16At(x*2+0, y*2+1),
				)
				t.Logf("(%d,%d): %x, %x", x*x, 2*y,
					src.Gray16At(x*2+1, y*2+1),
					src.Gray16At(x*2+1, y*2+0),
				)
				t.Fatalf("(%d,%d): %x != %x", x, y, v1, v2)
			}
		}
	}
}

func TestPyrDown_rgba(t *testing.T) {
	src := tToRGBA(tLoadImage("./testdata/lena.png"))

	dst0 := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx()/2, src.Bounds().Dy()/2))
	dst1 := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx()/2, src.Bounds().Dy()/2))

	draw.ApproxBiLinear.Scale(
		dst0, dst0.Bounds(),
		src, src.Bounds(),
		draw.Src, nil,
	)
	_PyrDown_ApproxBiLinear_RGBA_RGBA(
		dst1, dst1.Bounds(),
		src, image.Pt(0, 0),
	)

	for y := 0; y < src.Bounds().Dy()/2; y++ {
		for x := 0; x < src.Bounds().Dx()/2; x++ {
			if v1, v2 := dst0.RGBAAt(x, y), dst1.RGBAAt(x, y); v1 != v2 {
				t.Logf("(%d,%d): %x, %x", x*x, 2*y,
					src.RGBAAt(x*2+0, y*2+0),
					src.RGBAAt(x*2+0, y*2+1),
				)
				t.Logf("(%d,%d): %x, %x", x*x, 2*y,
					src.RGBAAt(x*2+1, y*2+1),
					src.RGBAAt(x*2+1, y*2+0),
				)
				t.Fatalf("(%d,%d): %x != %x", x, y, v1, v2)
			}
		}
	}
}

func TestPyrDown_rgba64(t *testing.T) {
	src := tToRGBA64(tLoadImage("./testdata/lena.png"))

	dst0 := image.NewRGBA64(image.Rect(0, 0, src.Bounds().Dx()/2, src.Bounds().Dy()/2))
	dst1 := image.NewRGBA64(image.Rect(0, 0, src.Bounds().Dx()/2, src.Bounds().Dy()/2))

	draw.ApproxBiLinear.Scale(
		dst0, dst0.Bounds(),
		src, src.Bounds(),
		draw.Src, nil,
	)
	_PyrDown_ApproxBiLinear_RGBA64_RGBA64(
		dst1, dst1.Bounds(),
		src, image.Pt(0, 0),
	)

	for y := 0; y < src.Bounds().Dy()/2; y++ {
		for x := 0; x < src.Bounds().Dx()/2; x++ {
			if v1, v2 := dst0.RGBA64At(x, y), dst1.RGBA64At(x, y); v1 != v2 {
				t.Logf("(%d,%d): %x, %x", x*x, 2*y,
					src.RGBA64At(x*2+0, y*2+0),
					src.RGBA64At(x*2+0, y*2+1),
				)
				t.Logf("(%d,%d): %x, %x", x*x, 2*y,
					src.RGBA64At(x*2+1, y*2+1),
					src.RGBA64At(x*2+1, y*2+0),
				)
				t.Fatalf("(%d,%d): %x != %x", x, y, v1, v2)
			}
		}
	}
}

func tLoadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("tLoadImage: os.Open(%q), err= %v", filename, err)
	}
	defer f.Close()

	m, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("tLoadImage: image.Decode, err= %v", err)
	}
	return m
}

func tToGray(m image.Image) *image.Gray {
	if p, ok := m.(*image.Gray); ok {
		return p
	}
	p := image.NewGray(m.Bounds())
	draw.Draw(p, p.Bounds(), m, image.Pt(0, 0), draw.Src)
	return p
}

func tToGray16(m image.Image) *image.Gray16 {
	if p, ok := m.(*image.Gray16); ok {
		return p
	}
	p := image.NewGray16(m.Bounds())
	draw.Draw(p, p.Bounds(), m, image.Pt(0, 0), draw.Src)
	return p
}

func tToRGBA(m image.Image) *image.RGBA {
	if p, ok := m.(*image.RGBA); ok {
		return p
	}
	p := image.NewRGBA(m.Bounds())
	draw.Draw(p, p.Bounds(), m, image.Pt(0, 0), draw.Src)
	return p
}

func tToRGBA64(m image.Image) *image.RGBA64 {
	if p, ok := m.(*image.RGBA64); ok {
		return p
	}
	p := image.NewRGBA64(m.Bounds())
	draw.Draw(p, p.Bounds(), m, image.Pt(0, 0), draw.Src)
	return p
}
