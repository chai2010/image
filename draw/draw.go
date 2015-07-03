// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"
	"image/color"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

type Image interface {
	image.Image
	Set(x, y int, c color.Color)
}

type Drawer interface {
	Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point)
}

func Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point) {
	if r.In(dst.Bounds()) && image.Rect(sp.X, sp.Y, sp.X+r.Dx(), sp.Y+r.Dy()).In(src.Bounds()) {
		// fast case
	}
	xdraw.Draw(dst, r, src, sp, xdraw.Src)
}
