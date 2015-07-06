// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"image"
	"image/draw"
	"reflect"

	xdraw "github.com/chai2010/image/draw"
)

type PyramidDriver interface {
	Channels() int
	DataType() reflect.Kind
	GetTile(level, col, row int) (draw.Image, error)
	SetTile(level, col, row int, tile image.Image) error
	xdraw.Drawer
	xdraw.Scaler
}
