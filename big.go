// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package memp

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	xdraw "golang.org/x/image/draw"
)

type SubImager interface {
	draw.Image
	SubImage(r image.Rectangle) image.Image
}

type BigImage struct {
	PyrDowner    PyrDowner
	Rect         image.Rectangle
	TileSize     image.Point
	TileMap      [][][]draw.Image // m.TileMap[level][col][row]
	TileMapMutex sync.Mutex
}

func NewBigImage(r image.Rectangle, tileSize image.Point, op PyrDowner) *BigImage {
	makeImageTileMap := func(r image.Rectangle, tileSize image.Point) (tileMap [][][]draw.Image) {
		xLevels := 0
		for i := 0; ; i++ {
			if x := (tileSize.X << uint8(i)); x >= r.Dx() {
				xLevels = i + 1
				break
			}
		}
		yLevels := 0
		for i := 0; ; i++ {
			if y := (tileSize.Y << uint8(i)); y >= r.Dy() {
				yLevels = i + 1
				break
			}
		}
		tileMap = make([][][]draw.Image, maxInt(xLevels, yLevels))
		for i := 0; i < len(tileMap); i++ {
			xTileSize := tileSize.X << uint8(len(tileMap)-i-1)
			yTileSize := tileSize.Y << uint8(len(tileMap)-i-1)
			xTilesNum := (r.Dx() + xTileSize - 1) / xTileSize
			yTilesNum := (r.Dy() + yTileSize - 1) / yTileSize

			tileMap[i] = make([][]draw.Image, xTilesNum)
			for x := 0; x < xTilesNum; x++ {
				tileMap[i][x] = make([]draw.Image, yTilesNum)
			}
		}
		return
	}

	m := &BigImage{
		PyrDowner: op,
		Rect:      r,
		TileMap:   makeImageTileMap(r, tileSize),
		TileSize:  tileSize,
	}
	return m
}

func (p *BigImage) SubLevels(levels int) *BigImage {
	r := p.Rect
	for i := levels; i < p.Levels(); i++ {
		r.Min.X /= 2
		r.Min.Y /= 2
		r.Max.X /= 2
		r.Max.Y /= 2
	}
	return &BigImage{
		PyrDowner:    p.PyrDowner,
		Rect:         r,
		TileMap:      p.TileMap[:levels],
		TileSize:     p.TileSize,
		TileMapMutex: p.TileMapMutex,
	}
}

func (p *BigImage) ColorModel() color.Model {
	return color.ModelFunc(func(c color.Color) color.Color {
		return c
	})
}

func (p *BigImage) Bounds() image.Rectangle { return p.Rect }

func (p *BigImage) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Gray{}
	}
	m := p.GetTile(p.Levels()-1, x/p.TileSize.X, y/p.TileSize.Y)
	c := m.At(x%p.TileSize.X, y%p.TileSize.Y)
	return c
}

func (p *BigImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	level, col, row := p.Levels()-1, x/p.TileSize.X, y/p.TileSize.Y
	p.GetTile(level, col, row).Set(x%p.TileSize.X, y%p.TileSize.Y, c)
	if p.TileEdgeOverlap {
		if x > 0 && y > 0 && x%p.TileSize.X == 0 && y%p.TileSize.Y == 0 {
			p.GetTile(level, col-1, row-1).Set(p.TileSize.X, p.TileSize.Y, c)
		}
		if x > 0 && x%p.TileSize.X == 0 {
			p.GetTile(level, col-1, row).Set(p.TileSize.X, y%p.TileSize.Y, c)
		}
		if y > 0 && y%p.TileSize.Y == 0 {
			p.GetTile(level, col, row-1).Set(x%p.TileSize.X, p.TileSize.Y, c)
		}
	}
	return
}

func (p *BigImage) Levels() int {
	return len(p.TileMap)
}

func (p *BigImage) adjustLevel(level int) int {
	if level < 0 {
		return p.Levels() + level
	}
	return level
}

func (p *BigImage) TilesAcross(level int) int {
	level = p.adjustLevel(level)
	v := len(p.TileMap[level])
	return v
}

func (p *BigImage) TilesDown(level int) int {
	level = p.adjustLevel(level)
	v := len(p.TileMap[level][0])
	return v
}

func (p *BigImage) GetTile(level, col, row int) (m draw.BigImage) {
	p.TileMapMutex.Lock()
	defer p.TileMapMutex.Unlock()
	level = p.adjustLevel(level)
	if m = p.TileMap[level][col][row]; m != nil {
		return
	}
	m = newImageTile(p.TileSize, p.Model, p.ZeroColor, p.TileEdgeOverlap)
	p.TileMap[level][col][row] = m
	return
}

func (p *BigImage) SetTile(level, col, row int, m draw.BigImage) (err error) {
	p.TileMapMutex.Lock()
	defer p.TileMapMutex.Unlock()

	level = p.adjustLevel(level)
	p.TileMap[level][col][row] = m
	return
}

func (p *BigImage) ReadRect(level int, r image.Rectangle, buf SubImager) (m image.BigImage, err error) {
	level = p.adjustLevel(level)

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

	if buf == nil {
		buf = newImageTile(r.Size(), p.Model, p.ZeroColor, p.TileEdgeOverlap)
	}

	var wg sync.WaitGroup
	for col := tMinCol; col < tMaxCol; col++ {
		for row := tMinRow; row < tMaxRow; row++ {
			wg.Add(1)
			go func(level, col, row int) {
				p.readRectFromTile(buf, p.GetTile(level, col, row), r.Min.X, r.Min.Y, r.Dx(), r.Dy(), col, row)
				wg.Done()
			}(level, col, row)
		}
	}
	wg.Wait()
	m = buf.SubImage(r)
	return
}

func (p *BigImage) readRectFromTile(dst, tile draw.BigImage, x, y, dx, dy, col, row int) {
	bMinX := x
	bMinY := y
	bMaxX := x + dx
	bMaxY := y + dy

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
		dst.(draw.BigImage), image.Rect(
			zMinX-bMinX,
			zMinY-bMinY,
			zMaxX-bMinX,
			zMaxY-bMinY,
		),
		tile, image.Pt(
			zMinX-tMinX,
			zMinY-tMinY,
		),
		xdraw.Src,
	)
	return
}

func (p *BigImage) WriteRect(level int, r image.Rectangle, m image.BigImage) (err error) {
	level = p.adjustLevel(level)
	r = r.Intersect(p.Bounds())
	if level < 0 || level >= p.Levels() {
		err = fmt.Errorf("image/big: BigImage.WriteRect, level = %v", level)
		return
	}
	if r.Empty() {
		return
	}

	tMinCol := r.Min.X / p.TileSize.X
	tMinRow := r.Min.Y / p.TileSize.Y
	tMaxCol := (r.Max.X + p.TileSize.X - 1) / p.TileSize.X
	tMaxRow := (r.Max.Y + p.TileSize.Y - 1) / p.TileSize.Y

	if p.TileEdgeOverlap {
		if r.Min.X > 0 && r.Min.X%p.TileSize.X == 0 {
			tMinCol--
		}
		if r.Min.Y > 0 && r.Min.Y%p.TileSize.Y == 0 {
			tMinRow--
		}
	}
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
				p.writeRectToTile(p.GetTile(level, col, row), m, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), col, row)
				wg.Done()
			}(level, col, row)
		}
	}
	wg.Wait()

	p.updateRectPyramid(level, r.Min.X, r.Min.Y, r.Dx(), r.Dy())
	return
}

func (p *BigImage) writeRectToTile(tile draw.BigImage, src image.BigImage, x, y, dx, dy, col, row int) {
	bMinX := x
	bMinY := y
	bMaxX := x + dx
	bMaxY := y + dy

	tMinX := col * p.TileSize.X
	tMinY := row * p.TileSize.Y
	tMaxX := tMinX + p.TileSize.X
	tMaxY := tMinY + p.TileSize.Y

	if p.TileEdgeOverlap {
		tMaxX++
		tMaxY++
	}

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
		xdraw.Src,
	)
	return
}

func (p *BigImage) updateRectPyramid(level, x, y, dx, dy int) {
	for level > 0 && dx > 0 && dy > 0 {
		minX, minY := x, y
		maxX, maxY := x+dx, y+dy

		tMinCol := minX / p.TileSize.X
		tMinRow := minY / p.TileSize.Y
		tMaxCol := (maxX + p.TileSize.X - 1) / p.TileSize.X
		tMaxRow := (maxY + p.TileSize.Y - 1) / p.TileSize.Y

		if p.TileEdgeOverlap {
			if minX > 0 && minX%p.TileSize.X == 0 {
				tMinCol--
			}
			if minY > 0 && minY%p.TileSize.Y == 0 {
				tMinRow--
			}
		}
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

func (p *BigImage) updateParentTile(level, col, row int) {
	switch {
	case col%2 == 0 && row%2 == 0:
		p.PyrDowner(
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
		p.PyrDowner(
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
		p.PyrDowner(
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
		p.PyrDowner(
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

func (p *BigImage) SelectTileList(level int, r image.Rectangle) (
	levelList, colList, rowList []int,
) {
	level = p.adjustLevel(level)
	if level < 0 || level >= p.Levels() {
		return
	}

	r = r.Intersect(p.Bounds())
	x, y, dx, dy := r.Min.X, r.Min.Y, r.Dx(), r.Dy()

	for level >= 0 && dx >= 0 && dy >= 0 {
		minX, minY := x, y
		maxX, maxY := x+dx, y+dy

		tMinCol := minX / p.TileSize.X
		tMinRow := minY / p.TileSize.Y
		tMaxCol := (maxX + p.TileSize.X - 1) / p.TileSize.X
		tMaxRow := (maxY + p.TileSize.Y - 1) / p.TileSize.Y

		if p.TileEdgeOverlap {
			if minX > 0 && minX%p.TileSize.X == 0 {
				tMinCol--
			}
			if minY > 0 && minY%p.TileSize.Y == 0 {
				tMinRow--
			}
		}
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
