// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pyramid

import (
	"image"
)

var (
	_ SubImager = (*TileImage)(nil)
)

type SubImager interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

type TileImage struct {
	ImageSize image.Point
	TileSize  image.Point
	SubImager
}

func NewTileImage(m SubImager, tileSize image.Point) *TileImage {
	return &TileImage{
		ImageSize: image.Pt(m.Bounds().Dx(), m.Bounds().Dy()),
		TileSize:  tileSize,
		SubImager: m,
	}
}

func (p *TileImage) TilesAcross() int {
	return (p.ImageSize.X + p.TileSize.X - 1) / p.TileSize.X
}

func (p *TileImage) TilesDown() int {
	return (p.ImageSize.Y + p.TileSize.Y - 1) / p.TileSize.Y
}

func (p *TileImage) GetTile(col, row int) (m image.Image) {
	minX := col * p.TileSize.X
	minY := row * p.TileSize.Y
	maxX := minX + p.TileSize.X
	maxY := minY + p.TileSize.Y
	r := image.Rect(minX, minY, maxX, maxY)
	return p.SubImager.SubImage(r)
}
