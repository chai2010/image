// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package draw

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
)

var (
	_ image.Image = (*Pyramid)(nil)
)

type Pyramid struct {
	ImageSize []image.Point // [0] is top
	TileSize  image.Point
	Driver    PyramidDriver
}

type PyramidDriver interface {
	ColorModel() color.Model
	GetTile(level, col, row int) (draw.Image, error)
	SetTile(level, col, row int, tile image.Image) error
	Scaler
}

func NewPyramid(imageSize, tileSize image.Point, driver PyramidDriver) *Pyramid {
	if v := imageSize; v.X <= 0 || v.Y <= 0 {
		panicf("draw: NewPyramid, imageSize = %v", imageSize)
	}
	if v := tileSize; v.X <= 0 || v.Y <= 0 {
		panicf("draw: NewPyramid, tileSize = %v", tileSize)
	}
	if driver == nil {
		panicf("draw: NewPyramid, driver = <nil>")
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

	p := &Pyramid{
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

func (p *Pyramid) SubLevels(levels int) *Pyramid {
	if levels <= 0 || levels > len(p.ImageSize) {
		panicf("draw: Pyramid.SubLevels, levels = %v", levels)
	}
	return &Pyramid{
		ImageSize: append([]image.Point{}, p.ImageSize[:levels]...),
		TileSize:  p.TileSize,
		Driver:    p.Driver,
	}
}

func (p *Pyramid) Bounds() image.Rectangle {
	return image.Rectangle{Max: p.ImageSize[len(p.ImageSize)-1]}
}

func (p *Pyramid) ColorModel() color.Model {
	return p.Driver.ColorModel()
}

func (p *Pyramid) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Bounds())) {
		return color.Gray{}
	}

	m, err := p.Driver.GetTile(len(p.ImageSize)-1, x/p.TileSize.X, y/p.TileSize.Y)
	if err != nil {
		panicf("draw: Pyramid.At(%d,%d), err = %v", x, y, err)
	}

	c := m.At(x%p.TileSize.X, y%p.TileSize.Y)
	return c
}

func (p *Pyramid) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Bounds())) {
		return
	}

	m, err := p.Driver.GetTile(len(p.ImageSize)-1, x/p.TileSize.X, y/p.TileSize.Y)
	if err != nil {
		panicf("draw: Pyramid.Set(%d,%d), err = %v", x, y, err)
	}

	m.Set(x%p.TileSize.X, y%p.TileSize.Y, c)
	return
}

func (p *Pyramid) Levels() int {
	return len(p.ImageSize)
}

func (p *Pyramid) adjustLevel(level int) int {
	if level < 0 {
		return len(p.ImageSize) + level
	}
	return level
}

func (p *Pyramid) TilesAcross(level int) int {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panicf("draw: Pyramid.TilesAcross, level = %v", level)
	}
	level = p.adjustLevel(level)
	v := (p.ImageSize[level].X + p.TileSize.X - 1) / p.TileSize.X
	return v
}

func (p *Pyramid) TilesDown(level int) int {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panicf("draw: Pyramid.TilesDown, level = %v", level)
	}
	level = p.adjustLevel(level)
	v := (p.ImageSize[level].Y + p.TileSize.Y - 1) / p.TileSize.Y
	return v
}

func (p *Pyramid) GetTile(level, col, row int) (m draw.Image) {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panicf("draw: Pyramid.GetTile, level = %v", level)
	}

	level = p.adjustLevel(level)
	if col < 0 || col >= (p.ImageSize[level].X+p.TileSize.X-1)/p.TileSize.X {
		panicf("draw: Pyramid.GetTile, level = %v, col = %v", level, col)
	}
	if row < 0 || row >= (p.ImageSize[level].X+p.TileSize.Y-1)/p.TileSize.Y {
		panicf("draw: Pyramid.GetTile, level = %v, row = %v", level, row)
	}

	m, err := p.Driver.GetTile(level, col, row)
	if err != nil {
		panicf("draw: Pyramid.GetTile(%d,%d,%d), err = %v", level, col, row, err)
	}
	return m
}

func (p *Pyramid) SetTile(level, col, row int, m image.Image) {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panic(errorf("draw: Pyramid.GetTile, level = %v", level))
	}

	level = p.adjustLevel(level)
	if col < 0 || col >= (p.ImageSize[level].X+p.TileSize.X-1)/p.TileSize.X {
		panic(errorf("draw: Pyramid.GetTile, level = %v, col = %v", level, col))
	}
	if row < 0 || row >= (p.ImageSize[level].X+p.TileSize.Y-1)/p.TileSize.Y {
		panic(errorf("draw: Pyramid.GetTile, level = %v, row = %v", level, row))
	}

	if b := m.Bounds(); b.Dx() != p.TileSize.X || b.Dy() != p.TileSize.Y {
		panic(errorf("draw: Pyramid.GetTile, m.Bounds() = %v", m.Bounds()))
	}

	err := p.Driver.SetTile(level, col, row, m)
	if err != nil {
		panicf("draw: Pyramid.SetTile(%d,%d,%d), err = %v", level, col, row, err)
	}
}

func (p *Pyramid) ReadRect(dst draw.Image, r image.Rectangle, level int) (err error) {
	level = p.adjustLevel(level)
	r = r.Intersect(dst.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		return errorf("draw: Pyramid.ReadRect, level = %v", level)
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

func (p *Pyramid) readRectFromTile(dst, tile draw.Image, r image.Rectangle, col, row int) {
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

	Draw(
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

func (p *Pyramid) WriteRect(m image.Image, r image.Rectangle, level int) (err error) {
	level = p.adjustLevel(level)
	r = r.Intersect(p.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		err = errorf("draw: Pyramid.WriteRect, level = %v", level)
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

func (p *Pyramid) writeRectToTile(tile draw.Image, src image.Image, r image.Rectangle, col, row int) {
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

	Draw(
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

func (p *Pyramid) updateRectPyramid(level, x, y, dx, dy int) {
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

func (p *Pyramid) updateParentTile(level, col, row int) {
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

func (p *Pyramid) SelectTileList(level int, r image.Rectangle) (
	levelList, colList, rowList []int,
) {
	level = p.adjustLevel(level)
	r = r.Intersect(p.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panic(errorf("draw: Pyramid.SelectTileList, level = %v", level))
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
