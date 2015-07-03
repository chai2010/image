// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	memp "github.com/chai2010/image"
	xdraw "github.com/chai2010/image/draw"
)

type Image struct {
	ImageSize []image.Point // [0] is top
	TileSize  image.Point
	Driver    Driver
}

func NewImage(imageSize image.Point, tileSize image.Point, driver Driver) *Image {
	if v := imageSize; v.X <= 0 || v.Y <= 0 {
		panicf("big: NewImage, imageSize = %v", imageSize)
	}
	if v := tileSize; v.X <= 0 || v.Y <= 0 {
		panicf("big: NewImage, tileSize = %v", tileSize)
	}
	if driver == nil {
		panicf("big: NewImage, driver = <nil>")
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

	p := &Image{
		ImageSize: make([]image.Point, maxInt(xLevels, yLevels)),
		TileSize:  tileSize,
		Driver:    driver,
	}
	for i, _ := range p.ImageSize {
		k := len(p.ImageSize) - i - 1
		p.ImageSize[k] = image.Point{
			X: imageSize.X >> uint8(i),
			Y: imageSize.Y >> uint8(i),
		}
		if p.ImageSize[k].X <= 0 {
			p.ImageSize[k].X = 1
		}
		if p.ImageSize[k].Y <= 0 {
			p.ImageSize[k].Y = 1
		}
	}

	return p
}

func (p *Image) SubLevels(levels int) *Image {
	if levels <= 0 || levels > len(p.ImageSize) {
		panicf("big: Image.SubLevels, levels = %v", levels)
	}
	return &Image{
		ImageSize: append([]image.Point{}, p.ImageSize[:levels]...),
		TileSize:  p.TileSize,
		Driver:    p.Driver,
	}
}

func (p *Image) Bounds() image.Rectangle {
	return image.Rectangle{Max: p.ImageSize[len(p.ImageSize)-1]}
}

func (p *Image) ColorModel() color.Model {
	return memp.ColorModel(p.Driver.Channels(), p.Driver.DataType())
}

func (p *Image) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Bounds())) {
		return color.Gray{}
	}

	m, err := p.Driver.GetTile(len(p.ImageSize)-1, x/p.TileSize.X, y/p.TileSize.Y)
	if err != nil {
		panicf("big: Image.At(%d,%d), err = %v", x, y, err)
	}

	c := m.At(x%p.TileSize.X, y%p.TileSize.Y)
	return c
}

func (p *Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Bounds())) {
		return
	}

	m, err := p.Driver.GetTile(len(p.ImageSize)-1, x/p.TileSize.X, y/p.TileSize.Y)
	if err != nil {
		panicf("big: Image.Set(%d,%d), err = %v", x, y, err)
	}

	m.Set(x%p.TileSize.X, y%p.TileSize.Y, c)
	return
}

func (p *Image) Levels() int {
	return len(p.ImageSize)
}

func (p *Image) adjustLevel(level int) int {
	if level < 0 {
		return len(p.ImageSize) + level
	}
	return level
}

func (p *Image) TilesAcross(level int) int {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panicf("big: Image.TilesAcross, level = %v", level)
	}
	level = p.adjustLevel(level)
	v := (p.ImageSize[level].X + p.TileSize.X - 1) / p.TileSize.X
	return v
}

func (p *Image) TilesDown(level int) int {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panicf("big: Image.TilesDown, level = %v", level)
	}
	level = p.adjustLevel(level)
	v := (p.ImageSize[level].Y + p.TileSize.Y - 1) / p.TileSize.Y
	return v
}

func (p *Image) GetTile(level, col, row int) (m draw.Image) {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panicf("big: Image.GetTile, level = %v", level)
	}

	level = p.adjustLevel(level)
	if col < 0 || col >= (p.ImageSize[level].X+p.TileSize.X-1)/p.TileSize.X {
		panicf("big: Image.GetTile, level = %v, col = %v", level, col)
	}
	if row < 0 || row >= (p.ImageSize[level].X+p.TileSize.Y-1)/p.TileSize.Y {
		panicf("big: Image.GetTile, level = %v, row = %v", level, row)
	}

	m, err := p.Driver.GetTile(level, col, row)
	if err != nil {
		panicf("big: Image.GetTile(%d,%d,%d), err = %v", level, col, row, err)
	}
	return m
}

func (p *Image) SetTile(level, col, row int, m image.Image) {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panic(fmt.Errorf("bigimg: Image.GetTile, level = %v", level))
	}

	level = p.adjustLevel(level)
	if col < 0 || col >= (p.ImageSize[level].X+p.TileSize.X-1)/p.TileSize.X {
		panic(fmt.Errorf("bigimg: Image.GetTile, level = %v, col = %v", level, col))
	}
	if row < 0 || row >= (p.ImageSize[level].X+p.TileSize.Y-1)/p.TileSize.Y {
		panic(fmt.Errorf("bigimg: Image.GetTile, level = %v, row = %v", level, row))
	}

	if b := m.Bounds(); b.Dx() != p.TileSize.X || b.Dy() != p.TileSize.Y {
		panic(fmt.Errorf("bigimg: Image.GetTile, m.Bounds() = %v", m.Bounds()))
	}

	err := p.Driver.SetTile(level, col, row, m)
	if err != nil {
		panicf("big: Image.SetTile(%d,%d,%d), err = %v", level, col, row, err)
	}
}

func (p *Image) ReadRect(dst draw.Image, r image.Rectangle, level int) (err error) {
	level = p.adjustLevel(level)
	r = r.Intersect(dst.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		return fmt.Errorf("bigimg: Image.ReadRect, level = %v", level)
	}
	if r.Empty() {
		return nil
	}

	tMinCol := r.Min.X / p.TileSize.X
	tMinRow := r.Min.Y / p.TileSize.Y
	tMaxCol := (r.Max.X + p.TileSize.X - 1) / p.TileSize.X
	tMaxRow := (r.Max.Y + p.TileSize.Y - 1) / p.TileSize.Y

	if max := p.TilesAcross(level); tMaxCol > max {
		tMaxCol = max
	}
	if max := p.TilesDown(level); tMaxRow > max {
		tMaxRow = max
	}

	var wg sync.WaitGroup
	for col := tMinCol; col < tMaxCol; col++ {
		for row := tMinRow; row < tMaxRow; row++ {
			wg.Add(1)
			go func(level, col, row int) {
				p.readRectFromTile(dst, p.GetTile(level, col, row), r, col, row)
				wg.Done()
			}(level, col, row)
		}
	}
	wg.Wait()
	return
}

func (p *Image) readRectFromTile(dst, tile draw.Image, r image.Rectangle, col, row int) {
	bMinX := r.Min.X
	bMinY := r.Min.Y
	bMaxX := r.Max.X
	bMaxY := r.Max.Y

	tMinX := col * p.TileSize.X
	tMinY := row * p.TileSize.Y
	tMaxX := tMinX + p.TileSize.X
	tMaxY := tMinY + p.TileSize.Y

	zMinX := maxInt(bMinX, tMinX)
	zMinY := maxInt(bMinY, tMinY)
	zMaxX := minInt(bMaxX, tMaxX)
	zMaxY := minInt(bMaxY, tMaxY)

	if zMinX >= zMaxX || zMinY >= zMaxY {
		return
	}

	xdraw.Draw(
		dst, image.Rect(
			zMinX-bMinX,
			zMinY-bMinY,
			zMaxX-bMinX,
			zMaxY-bMinY,
		),
		tile, image.Pt(
			zMinX-tMinX,
			zMinY-tMinY,
		),
	)
	return
}

func (p *Image) WriteRect(m image.Image, r image.Rectangle, level int) (err error) {
	level = p.adjustLevel(level)
	r = r.Intersect(p.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		err = fmt.Errorf("bigimg: Image.WriteRect, level = %v", level)
		return
	}
	if r.Empty() {
		return
	}

	tMinCol := r.Min.X / p.TileSize.X
	tMinRow := r.Min.Y / p.TileSize.Y
	tMaxCol := (r.Max.X + p.TileSize.X - 1) / p.TileSize.X
	tMaxRow := (r.Max.Y + p.TileSize.Y - 1) / p.TileSize.Y

	if max := p.TilesAcross(level); tMaxCol > max {
		tMaxCol = max
	}
	if max := p.TilesDown(level); tMaxRow > max {
		tMaxRow = max
	}

	var wg sync.WaitGroup
	for col := tMinCol; col < tMaxCol; col++ {
		for row := tMinRow; row < tMaxRow; row++ {
			wg.Add(1)
			go func(level, col, row int) {
				p.writeRectToTile(p.GetTile(level, col, row), m, r, col, row)
				wg.Done()
			}(level, col, row)
		}
	}
	wg.Wait()

	p.updateRectPyramid(level, r.Min.X, r.Min.Y, r.Dx(), r.Dy())
	return
}

func (p *Image) writeRectToTile(tile draw.Image, src image.Image, r image.Rectangle, col, row int) {
	bMinX := r.Min.X
	bMinY := r.Min.Y
	bMaxX := r.Max.X
	bMaxY := r.Max.Y

	tMinX := col * p.TileSize.X
	tMinY := row * p.TileSize.Y
	tMaxX := tMinX + p.TileSize.X
	tMaxY := tMinY + p.TileSize.Y

	zMinX := maxInt(bMinX, tMinX)
	zMinY := maxInt(bMinY, tMinY)
	zMaxX := minInt(bMaxX, tMaxX)
	zMaxY := minInt(bMaxY, tMaxY)

	if zMinX >= zMaxX || zMinY >= zMaxY {
		return
	}

	xdraw.Draw(
		tile, image.Rect(
			zMinX-tMinX,
			zMinY-tMinY,
			zMaxX-tMinX,
			zMaxY-tMinY,
		),
		src, image.Pt(
			zMinX-bMinX,
			zMinY-bMinY,
		),
	)
	return
}

func (p *Image) updateRectPyramid(level, x, y, dx, dy int) {
	for level > 0 && dx > 0 && dy > 0 {
		minX, minY := x, y
		maxX, maxY := x+dx, y+dy

		tMinCol := minX / p.TileSize.X
		tMinRow := minY / p.TileSize.Y
		tMaxCol := (maxX + p.TileSize.X - 1) / p.TileSize.X
		tMaxRow := (maxY + p.TileSize.Y - 1) / p.TileSize.Y

		if max := p.TilesAcross(level); tMaxCol > max {
			tMaxCol = max
		}
		if max := p.TilesDown(level); tMaxRow > max {
			tMaxRow = max
		}

		var wg sync.WaitGroup
		for row := tMinRow; row < tMaxRow; row++ {
			for col := tMinCol; col < tMaxCol; col++ {
				wg.Add(1)
				go func(level, col, row int) {
					p.updateParentTile(level, col, row)
					wg.Done()
				}(level, col, row)
			}
		}
		wg.Wait()

		x, dx = (minX+1)/2, (maxX-minX+1)/2
		y, dy = (minY+1)/2, (maxY-minY+1)/2
		level--
	}
	return
}

func (p *Image) updateParentTile(level, col, row int) {
	sr := image.Rect(0, 0, p.TileSize.X, p.TileSize.Y)
	switch {
	case col%2 == 0 && row%2 == 0:
		p.Driver.Scale(
			p.GetTile(level-1, col/2, row/2),
			image.Rect(
				(p.TileSize.X/2)*0,
				(p.TileSize.Y/2)*0,
				(p.TileSize.X/2)*0+p.TileSize.X/2,
				(p.TileSize.Y/2)*0+p.TileSize.Y/2,
			),
			p.GetTile(level, col, row),
			sr,
		)
	case col%2 == 0 && row%2 == 1:
		p.Driver.Scale(
			p.GetTile(level-1, col/2, row/2),
			image.Rect(
				(p.TileSize.X/2)*0,
				(p.TileSize.Y/2)*1,
				(p.TileSize.X/2)*0+p.TileSize.X/2,
				(p.TileSize.Y/2)*1+p.TileSize.Y/2,
			),
			p.GetTile(level, col, row),
			sr,
		)
	case col%2 == 1 && row%2 == 1:
		p.Driver.Scale(
			p.GetTile(level-1, col/2, row/2),
			image.Rect(
				(p.TileSize.X/2)*1,
				(p.TileSize.Y/2)*1,
				(p.TileSize.X/2)*1+p.TileSize.X/2,
				(p.TileSize.Y/2)*1+p.TileSize.Y/2,
			),
			p.GetTile(level, col, row),
			sr,
		)
	case col%2 == 1 && row%2 == 0:
		p.Driver.Scale(
			p.GetTile(level-1, col/2, row/2),
			image.Rect(
				(p.TileSize.X/2)*1,
				(p.TileSize.Y/2)*0,
				(p.TileSize.X/2)*1+p.TileSize.X/2,
				(p.TileSize.Y/2)*0+p.TileSize.Y/2,
			),
			p.GetTile(level, col, row),
			sr,
		)
	}
}

func (p *Image) SelectTileList(level int, r image.Rectangle) (
	levelList, colList, rowList []int,
) {
	level = p.adjustLevel(level)
	r = r.Intersect(p.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panic(fmt.Errorf("bigimg: Image.SelectTileList, level = %v", level))
	}
	if r.Empty() {
		return
	}

	x, y, dx, dy := r.Min.X, r.Min.Y, r.Dx(), r.Dy()

	for level >= 0 && dx >= 0 && dy >= 0 {
		minX, minY := x, y
		maxX, maxY := x+dx, y+dy

		tMinCol := minX / p.TileSize.X
		tMinRow := minY / p.TileSize.Y
		tMaxCol := (maxX + p.TileSize.X - 1) / p.TileSize.X
		tMaxRow := (maxY + p.TileSize.Y - 1) / p.TileSize.Y

		if max := p.TilesAcross(level); tMaxCol > max {
			tMaxCol = max
		}
		if max := p.TilesDown(level); tMaxRow > max {
			tMaxRow = max
		}

		for row := tMinRow; row <= tMaxRow; row++ {
			for col := tMinCol; col <= tMaxCol; col++ {
				levelList = append(levelList, level)
				colList = append(colList, col)
				rowList = append(rowList, row)
			}
		}

		x, dx = (minX+1)/2, (maxX-minX+1)/2
		y, dy = (minY+1)/2, (maxY-minY+1)/2
		level--
	}
	return
}
