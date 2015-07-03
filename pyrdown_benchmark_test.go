// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memp

import (
	"image"
	"testing"

	"golang.org/x/image/draw"
)

func BenchmarkPyrDownGray_16x16(b *testing.B)          { benchPyrDownGray(b, 16, 16) }
func BenchmarkPyrDownGray_16x16_x_draw(b *testing.B)   { benchPyrDownGray_x_draw(b, 16, 16) }
func BenchmarkPyrDownGray_32x32(b *testing.B)          { benchPyrDownGray(b, 32, 32) }
func BenchmarkPyrDownGray_32x32_x_draw(b *testing.B)   { benchPyrDownGray_x_draw(b, 32, 32) }
func BenchmarkPyrDownGray_64x64(b *testing.B)          { benchPyrDownGray(b, 64, 64) }
func BenchmarkPyrDownGray_64x64_x_draw(b *testing.B)   { benchPyrDownGray_x_draw(b, 64, 64) }
func BenchmarkPyrDownGray_128x128(b *testing.B)        { benchPyrDownGray(b, 128, 128) }
func BenchmarkPyrDownGray_128x128_x_draw(b *testing.B) { benchPyrDownGray_x_draw(b, 128, 128) }
func BenchmarkPyrDownGray_256x256(b *testing.B)        { benchPyrDownGray(b, 256, 256) }
func BenchmarkPyrDownGray_256x256_x_draw(b *testing.B) { benchPyrDownGray_x_draw(b, 256, 256) }

func BenchmarkPyrDownRGBA_16x16(b *testing.B)          { benchPyrDownRGBA(b, 16, 16) }
func BenchmarkPyrDownRGBA_16x16_x_draw(b *testing.B)   { benchPyrDownRGBA_x_draw(b, 16, 16) }
func BenchmarkPyrDownRGBA_32x32(b *testing.B)          { benchPyrDownRGBA(b, 32, 32) }
func BenchmarkPyrDownRGBA_32x32_x_draw(b *testing.B)   { benchPyrDownRGBA_x_draw(b, 32, 32) }
func BenchmarkPyrDownRGBA_64x64(b *testing.B)          { benchPyrDownRGBA(b, 64, 64) }
func BenchmarkPyrDownRGBA_64x64_x_draw(b *testing.B)   { benchPyrDownRGBA_x_draw(b, 64, 64) }
func BenchmarkPyrDownRGBA_128x128(b *testing.B)        { benchPyrDownRGBA(b, 128, 128) }
func BenchmarkPyrDownRGBA_128x128_x_draw(b *testing.B) { benchPyrDownRGBA_x_draw(b, 128, 128) }
func BenchmarkPyrDownRGBA_256x256(b *testing.B)        { benchPyrDownRGBA(b, 256, 256) }
func BenchmarkPyrDownRGBA_256x256_x_draw(b *testing.B) { benchPyrDownRGBA_x_draw(b, 256, 256) }

func BenchmarkPyrDownGray16_16x16(b *testing.B)          { benchPyrDownGray16(b, 16, 16) }
func BenchmarkPyrDownGray16_16x16_x_draw(b *testing.B)   { benchPyrDownGray16_x_draw(b, 16, 16) }
func BenchmarkPyrDownGray16_32x32(b *testing.B)          { benchPyrDownGray16(b, 32, 32) }
func BenchmarkPyrDownGray16_32x32_x_draw(b *testing.B)   { benchPyrDownGray16_x_draw(b, 32, 32) }
func BenchmarkPyrDownGray16_64x64(b *testing.B)          { benchPyrDownGray16(b, 64, 64) }
func BenchmarkPyrDownGray16_64x64_x_draw(b *testing.B)   { benchPyrDownGray16_x_draw(b, 64, 64) }
func BenchmarkPyrDownGray16_128x128(b *testing.B)        { benchPyrDownGray16(b, 128, 128) }
func BenchmarkPyrDownGray16_128x128_x_draw(b *testing.B) { benchPyrDownGray16_x_draw(b, 128, 128) }
func BenchmarkPyrDownGray16_256x256(b *testing.B)        { benchPyrDownGray16(b, 256, 256) }
func BenchmarkPyrDownGray16_256x256_x_draw(b *testing.B) { benchPyrDownGray16_x_draw(b, 256, 256) }

func BenchmarkPyrDownRGBA64_16x16(b *testing.B)          { benchPyrDownRGBA64(b, 16, 16) }
func BenchmarkPyrDownRGBA64_16x16_x_draw(b *testing.B)   { benchPyrDownRGBA64_x_draw(b, 16, 16) }
func BenchmarkPyrDownRGBA64_32x32(b *testing.B)          { benchPyrDownRGBA64(b, 32, 32) }
func BenchmarkPyrDownRGBA64_32x32_x_draw(b *testing.B)   { benchPyrDownRGBA64_x_draw(b, 32, 32) }
func BenchmarkPyrDownRGBA64_64x64(b *testing.B)          { benchPyrDownRGBA64(b, 64, 64) }
func BenchmarkPyrDownRGBA64_64x64_x_draw(b *testing.B)   { benchPyrDownRGBA64_x_draw(b, 64, 64) }
func BenchmarkPyrDownRGBA64_128x128(b *testing.B)        { benchPyrDownRGBA64(b, 128, 128) }
func BenchmarkPyrDownRGBA64_128x128_x_draw(b *testing.B) { benchPyrDownRGBA64_x_draw(b, 128, 128) }
func BenchmarkPyrDownRGBA64_256x256(b *testing.B)        { benchPyrDownRGBA64(b, 256, 256) }
func BenchmarkPyrDownRGBA64_256x256_x_draw(b *testing.B) { benchPyrDownRGBA64_x_draw(b, 256, 256) }

func benchPyrDownGray(b *testing.B, width, height int) {
	dst := image.NewGray(image.Rect(0, 0, width/2, height/2))
	src := image.NewGray(image.Rect(0, 0, width, height))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_PyrDown_ApproxBiLinear_Gray_Gray(
			dst, dst.Bounds(),
			src, image.Pt(0, 0),
		)
	}
}

func benchPyrDownGray_x_draw(b *testing.B, width, height int) {
	dst := image.NewGray(image.Rect(0, 0, width/2, height/2))
	src := image.NewGray(image.Rect(0, 0, width, height))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		draw.ApproxBiLinear.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Src, nil,
		)
	}
}

func benchPyrDownRGBA(b *testing.B, width, height int) {
	dst := image.NewRGBA(image.Rect(0, 0, width/2, height/2))
	src := image.NewRGBA(image.Rect(0, 0, width, height))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_PyrDown_ApproxBiLinear_RGBA_RGBA(
			dst, dst.Bounds(),
			src, image.Pt(0, 0),
		)
	}
}

func benchPyrDownRGBA_x_draw(b *testing.B, width, height int) {
	dst := image.NewRGBA(image.Rect(0, 0, width/2, height/2))
	src := image.NewRGBA(image.Rect(0, 0, width, height))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		draw.ApproxBiLinear.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Src, nil,
		)
	}
}

func benchPyrDownGray16(b *testing.B, width, height int) {
	dst := image.NewGray16(image.Rect(0, 0, width/2, height/2))
	src := image.NewGray16(image.Rect(0, 0, width, height))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_PyrDown_ApproxBiLinear_Gray16_Gray16(
			dst, dst.Bounds(),
			src, image.Pt(0, 0),
		)
	}
}

func benchPyrDownGray16_x_draw(b *testing.B, width, height int) {
	dst := image.NewGray16(image.Rect(0, 0, width/2, height/2))
	src := image.NewGray16(image.Rect(0, 0, width, height))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		draw.ApproxBiLinear.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Src, nil,
		)
	}
}

func benchPyrDownRGBA64(b *testing.B, width, height int) {
	dst := image.NewRGBA64(image.Rect(0, 0, width/2, height/2))
	src := image.NewRGBA64(image.Rect(0, 0, width, height))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_PyrDown_ApproxBiLinear_RGBA64_RGBA64(
			dst, dst.Bounds(),
			src, image.Pt(0, 0),
		)
	}
}

func benchPyrDownRGBA64_x_draw(b *testing.B, width, height int) {
	dst := image.NewRGBA64(image.Rect(0, 0, width/2, height/2))
	src := image.NewRGBA64(image.Rect(0, 0, width, height))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		draw.ApproxBiLinear.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Src, nil,
		)
	}
}
