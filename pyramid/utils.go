// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pyramid

import (
	"fmt"
	"image"
)

func PyramidLevels(imageSize image.Point, tileSize image.Point) int {
	if v := imageSize; v.X <= 0 || v.Y <= 0 {
		panic(fmt.Errorf("image/pyramid: PyramidLevels, imageSize = %v", imageSize))
	}
	if v := tileSize; v.X <= 0 || v.Y <= 0 {
		panic(fmt.Errorf("image/pyramid: PyramidLevels, tileSize = %v", tileSize))
	}

	xLevels := 0
	for i := 0; ; i++ {
		if x := (tileSize.X << uint8(i)); x >= imageSize.X {
			xLevels = i + 1
			break
		}
	}
	yLevels := 0
	for i := 0; ; i++ {
		if y := (tileSize.Y << uint8(i)); y >= imageSize.Y {
			yLevels = i + 1
			break
		}
	}

	return maxInt(xLevels, yLevels)
}

func minInt(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}
