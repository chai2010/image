// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"

	xdraw "github.com/chai2010/image/draw"
)

var (
	_ Driver     = (*Image)(nil)
	_ draw.Image = (*Image)(nil)
)

// Image is a simple big image base on the local file system.
type Image struct {
	filename string
	rect     image.Rectangle
	channels int
	dataType reflect.Kind
	pyramid  *Pyramid
}

func OpenImage(fs FileSystem, filename string) (p *Image, err error) {
	return
}

func CreateImage(fs FileSystem, filename string, imageSize image.Point, tileSize image.Point) (p *Image, err error) {
	return
}

func (p *Image) Close() error {
	panic("TODO")
}

func (p *Image) Channels() int {
	panic("TODO")
}
func (p *Image) DataType() reflect.Kind {
	panic("TODO")
}

func (p *Image) Levels() int {
	panic("TODO")
}

func (p *Image) SubLevels(levels int) *Pyramid {
	panic("TODO")
}

func (p *Image) TilesAcross(level int) int {
	panic("TODO")
}
func (p *Image) TilesDown(level int) int {
	panic("TODO")
}

func (p *Image) Bounds() image.Rectangle {
	panic("TODO")
}
func (p *Image) ColorModel() color.Model {
	panic("TODO")
}

func (p *Image) At(x, y int) color.Color {
	panic("TODO")
}
func (p *Image) Set(x, y int, c color.Color) {
	panic("TODO")
}

func (p *Image) SelectByRect(level int, r image.Rectangle) (levelList, colList, rowList []int) {
	panic("TODO")
}

func (p *Image) ReadRect(dst draw.Image, r image.Rectangle, level int) (err error) {
	panic("TODO")
}

func (p *Image) WriteRect(m image.Image, r image.Rectangle, level int) (err error) {
	panic("TODO")
}

func (p *Image) GetTile(level, col, row int) (draw.Image, error) {
	panic("TODO")
}
func (p *Image) SetTile(level, col, row int, tile image.Image) error {
	panic("TODO")
}
func (p *Image) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	xdraw.Draw(dst, r, src, sp)
}
func (p *Image) Scale(dst draw.Image, dr image.Rectangle, src image.Image, sr image.Rectangle) {
	xdraw.ApproxBiLinear.Scale(dst, dr, src, sr)
}

func (p *Image) Flush() error {
	panic("TODO")
}
