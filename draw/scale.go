// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

var (
	// NearestNeighbor is the nearest neighbor interpolator. It is very fast,
	// but usually gives very low quality results. When scaling up, the result
	// will look 'blocky'.
	NearestNeighbor = Scaler(nnScaler{})

	// ApproxBiLinear is a mixture of the nearest neighbor and bi-linear
	// interpolators. It is fast, but usually gives medium quality results.
	//
	// It implements bi-linear interpolation when upscaling and a bi-linear
	// blend of the 4 nearest neighbor pixels when downscaling. This yields
	// nicer quality than nearest neighbor interpolation when upscaling, but
	// the time taken is independent of the number of source pixels, unlike the
	// bi-linear interpolator. When downscaling a large image, the performance
	// difference can be significant.
	ApproxBiLinear = Scaler(abScaler{})
)

type Scaler interface {
	Scale(dst draw.Image, dr image.Rectangle, src image.Image, sr image.Rectangle)
}

type nnScaler struct{}

func (nnScaler) Scale(dst draw.Image, dr image.Rectangle, src image.Image, sr image.Rectangle) {
	if dr.In(dst.Bounds()) && sr.In(src.Bounds()) && sr.Dx() == dr.Dx()*2 && sr.Dy() == dr.Dy()*2 {
		nnPyrDownFast(dst, dr, src, sr.Min)
		return
	}
	xdraw.NearestNeighbor.Scale(
		dst, dr, src, sr,
		xdraw.Src, nil,
	)
}

type abScaler struct{}

func (abScaler) Scale(dst draw.Image, dr image.Rectangle, src image.Image, sr image.Rectangle) {
	if dr.In(dst.Bounds()) && sr.In(src.Bounds()) && sr.Dx() == dr.Dx()*2 && sr.Dy() == dr.Dy()*2 {
		abPyrDownFast(dst, dr, src, sr.Min)
		return
	}
	xdraw.ApproxBiLinear.Scale(
		dst, dr, src, sr,
		xdraw.Src, nil,
	)
}
