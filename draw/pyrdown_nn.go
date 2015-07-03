// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"

	"golang.org/x/image/draw"
)

func _PyrDown_NearestNeighbor(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	draw.NearestNeighbor.Scale(
		dst, r,
		src, image.Rect(sp.X, sp.Y, sp.X+r.Dx()/2, sp.Y+r.Dy()/2),
		draw.Src, nil,
	)
}
