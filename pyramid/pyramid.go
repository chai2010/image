// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pyramid provides tiled pyramid image support.
package pyramid

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	xdraw "github.com/chai2010/image/draw"
)

type Pyramid struct {
	ImageSize []image.Point // [0] is top
	TileSize  image.Point
	Driver    PyramidDriver
}

type PyramidDriver interface {
	ColorModel() color.Model
	GetTile(level, col, row int) draw.Image
	SetTile(level, col, row int, tile image.Image)
	PyrDown(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point)
}

func NewPyramid(imageSize image.Point, tileSize image.Point, driver PyramidDriver) *Pyramid {
	if v := imageSize; v.X <= 0 || v.Y <= 0 {
		panic(fmt.Errorf("image/pyramid: NewPyramid, imageSize = %v", imageSize))
	}
	if v := tileSize; v.X <= 0 || v.Y <= 0 {
		panic(fmt.Errorf("image/pyramid: NewPyramid, tileSize = %v", tileSize))
	}
	if driver == nil {
		panic(fmt.Errorf("image/pyramid: NewPyramid, driver = <nil>"))
	}

	p := &Pyramid{
		ImageSize: make([]image.Point, PyramidLevels(imageSize, tileSize)),
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
		panic(fmt.Errorf("image/pyramid: Image.SubLevels, levels = %v", levels))
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

	m := p.Driver.GetTile(len(p.ImageSize)-1, x/p.TileSize.X, y/p.TileSize.Y)
	if m == nil || m.Bounds().Empty() {
		panic(fmt.Errorf("image/pyramid: Image.At(%v,%v), p.Driver.GetTile return <nil>", x, y))
	}

	c := m.At(x%p.TileSize.X, y%p.TileSize.Y)
	return c
}

func (p *Pyramid) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Bounds())) {
		return
	}

	m := p.Driver.GetTile(len(p.ImageSize)-1, x/p.TileSize.X, y/p.TileSize.Y)
	if m == nil || m.Bounds().Empty() {
		panic(fmt.Errorf("image/pyramid: Image.Set(%v,%v,%v), p.Driver.GetTile return <nil>", x, y, c))
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
		panic(fmt.Errorf("image/pyramid: Image.TilesAcross, level = %v", level))
	}
	level = p.adjustLevel(level)
	v := (p.ImageSize[level].X + p.TileSize.X - 1) / p.TileSize.X
	return v
}

func (p *Pyramid) TilesDown(level int) int {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panic(fmt.Errorf("image/pyramid: Image.TilesDown, level = %v", level))
	}
	level = p.adjustLevel(level)
	v := (p.ImageSize[level].Y + p.TileSize.Y - 1) / p.TileSize.Y
	return v
}

func (p *Pyramid) GetTile(level, col, row int) (m draw.Image) {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panic(fmt.Errorf("image/pyramid: Image.GetTile, level = %v", level))
	}

	level = p.adjustLevel(level)
	if col < 0 || col >= (p.ImageSize[level].X+p.TileSize.X-1)/p.TileSize.X {
		panic(fmt.Errorf("image/pyramid: Image.GetTile, level = %v, col = %v", level, col))
	}
	if row < 0 || row >= (p.ImageSize[level].X+p.TileSize.Y-1)/p.TileSize.Y {
		panic(fmt.Errorf("image/pyramid: Image.GetTile, level = %v, row = %v", level, row))
	}

	return p.Driver.GetTile(level, col, row)
}

func (p *Pyramid) SetTile(level, col, row int, m image.Image) {
	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panic(fmt.Errorf("image/pyramid: Image.GetTile, level = %v", level))
	}

	level = p.adjustLevel(level)
	if col < 0 || col >= (p.ImageSize[level].X+p.TileSize.X-1)/p.TileSize.X {
		panic(fmt.Errorf("image/pyramid: Image.GetTile, level = %v, col = %v", level, col))
	}
	if row < 0 || row >= (p.ImageSize[level].X+p.TileSize.Y-1)/p.TileSize.Y {
		panic(fmt.Errorf("image/pyramid: Image.GetTile, level = %v, row = %v", level, row))
	}

	if b := m.Bounds(); b.Dx() != p.TileSize.X || b.Dy() != p.TileSize.Y {
		panic(fmt.Errorf("image/pyramid: Image.GetTile, m.Bounds() = %v", m.Bounds()))
	}

	p.Driver.SetTile(level, col, row, m)
}

func (p *Pyramid) ReadRect(dst draw.Image, r image.Rectangle, level int) (err error) {
	level = p.adjustLevel(level)
	r = r.Intersect(dst.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		return fmt.Errorf("image/pyramid: Image.ReadRect, level = %v", level)
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

func (p *Pyramid) WriteRect(m image.Image, r image.Rectangle, level int) (err error) {
	return p.writeRect(m, r, level, 32) // update all levels
}

func (p *Pyramid) writeRect(
	m image.Image, r image.Rectangle,
	level, updateLevelsLimit int,
) (err error) {
	level = p.adjustLevel(level)
	r = r.Intersect(p.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		err = fmt.Errorf("image/pyramid: Image.writeRect, level = %v", level)
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

	p.updateRectPyramid(level, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), updateLevelsLimit)
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

func (p *Pyramid) updateRectPyramid(level, x, y, dx, dy, updateLevelsLimit int) {
	for cnt := 0; cnt < updateLevelsLimit && level > 0 && dx > 0 && dy > 0; cnt++ {
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
		var workLimits = make(chan struct{}, 32)
		for row := tMinRow; row < tMaxRow; row++ {
			for col := tMinCol; col < tMaxCol; col++ {
				wg.Add(1)
				go func(level, col, row int) {
					workLimits <- struct{}{}
					defer func() { <-workLimits }()
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
	switch {
	case col%2 == 0 && row%2 == 0:
		p.Driver.PyrDown(
			p.GetTile(level-1, col/2, row/2),
			image.Rect(
				(p.TileSize.X/2)*0,
				(p.TileSize.Y/2)*0,
				(p.TileSize.X/2)*0+p.TileSize.X/2,
				(p.TileSize.Y/2)*0+p.TileSize.Y/2,
			),
			p.GetTile(level, col, row),
			image.Pt(0, 0),
		)
	case col%2 == 0 && row%2 == 1:
		p.Driver.PyrDown(
			p.GetTile(level-1, col/2, row/2),
			image.Rect(
				(p.TileSize.X/2)*0,
				(p.TileSize.Y/2)*1,
				(p.TileSize.X/2)*0+p.TileSize.X/2,
				(p.TileSize.Y/2)*1+p.TileSize.Y/2,
			),
			p.GetTile(level, col, row),
			image.Pt(0, 0),
		)
	case col%2 == 1 && row%2 == 1:
		p.Driver.PyrDown(
			p.GetTile(level-1, col/2, row/2),
			image.Rect(
				(p.TileSize.X/2)*1,
				(p.TileSize.Y/2)*1,
				(p.TileSize.X/2)*1+p.TileSize.X/2,
				(p.TileSize.Y/2)*1+p.TileSize.Y/2,
			),
			p.GetTile(level, col, row),
			image.Pt(0, 0),
		)
	case col%2 == 1 && row%2 == 0:
		p.Driver.PyrDown(
			p.GetTile(level-1, col/2, row/2),
			image.Rect(
				(p.TileSize.X/2)*1,
				(p.TileSize.Y/2)*0,
				(p.TileSize.X/2)*1+p.TileSize.X/2,
				(p.TileSize.Y/2)*0+p.TileSize.Y/2,
			),
			p.GetTile(level, col, row),
			image.Pt(0, 0),
		)
	}
}

func (p *Pyramid) SelectTileList(level int, r image.Rectangle) (
	levelList, colList, rowList []int,
) {
	level = p.adjustLevel(level)
	r = r.Intersect(p.Bounds())

	if level >= len(p.ImageSize) || level < -len(p.ImageSize) {
		panic(fmt.Errorf("image/pyramid: Image.SelectTileList, level = %v", level))
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

func (p *Pyramid) UpdatePyramid() (err error) {
	defer func() {
		if x := recover(); x != nil {
			if errx, ok := x.(error); ok {
				err = errx
				return
			} else {
				err = fmt.Errorf("image/pyramid: Pyramid.UpdatePyramid, %v", x)
				return
			}
		}
	}()

	var (
		stepX     = p.TileSize.X * 16
		stepY     = p.TileSize.Y * 16
		stepLevel = 3 // 16x16 => 4x4 => 1x1
	)
	for level := p.Levels() - 1; level > 0; level -= stepLevel {
		for x := 0; x < p.ImageSize[level].X; x += stepX {
			for y := 0; y < p.ImageSize[level].Y; y += stepY {
				p.updateRectPyramid(level, x, y, stepX, stepY, stepLevel)
			}
		}
	}
	return nil
}
