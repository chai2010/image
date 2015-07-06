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
	_ Driver     = (*File)(nil)
	_ draw.Image = (*File)(nil)
)

// File is a simple big image base on the local file system.
type File struct {
	filename string
	rect     image.Rectangle
	channels int
	dataType reflect.Kind
	pyramid  *Pyramid
}

func OpenFile(filename string) (f *File, err error) {
	return
}

func CreateFile(filename string, imageSize image.Point, tileSize image.Point) (f *File, err error) {
	return
}

func (f *File) Close() error {
	panic("TODO")
}

func (f *File) Channels() int {
	panic("TODO")
}
func (f *File) DataType() reflect.Kind {
	panic("TODO")
}

func (f *File) Levels() int {
	panic("TODO")
}

func (f *File) SubLevels(levels int) *Pyramid {
	panic("TODO")
}

func (f *File) TilesAcross(level int) int {
	panic("TODO")
}
func (f *File) TilesDown(level int) int {
	panic("TODO")
}

func (f *File) Bounds() image.Rectangle {
	panic("TODO")
}
func (f *File) ColorModel() color.Model {
	panic("TODO")
}

func (f *File) At(x, y int) color.Color {
	panic("TODO")
}
func (f *File) Set(x, y int, c color.Color) {
	panic("TODO")
}

func (f *File) SelectByRect(level int, r image.Rectangle) (levelList, colList, rowList []int) {
	panic("TODO")
}

func (f *File) ReadRect(dst draw.Image, r image.Rectangle, level int) (err error) {
	panic("TODO")
}

func (f *File) WriteRect(m image.Image, r image.Rectangle, level int) (err error) {
	panic("TODO")
}

func (f *File) GetTile(level, col, row int) (draw.Image, error) {
	panic("TODO")
}
func (f *File) SetTile(level, col, row int, tile image.Image) error {
	panic("TODO")
}
func (f *File) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	xdraw.Draw(dst, r, src, sp)
}
func (f *File) Scale(dst draw.Image, dr image.Rectangle, src image.Image, sr image.Rectangle) {
	xdraw.ApproxBiLinear.Scale(dst, dr, src, sr)
}

func (f *File) Flush() error {
	panic("TODO")
}
