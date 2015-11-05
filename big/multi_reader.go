// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"image"
	"reflect"

	ximage "github.com/chai2010/image"
	xdraw "github.com/chai2010/image/draw"
)

var (
	_ ImageReader = (*_MultiImageReader)(nil)
)

type _MultiImageReader struct {
	readers  map[image.Rectangle]ImageReader
	rect     image.Rectangle
	channels int
	dataType reflect.Kind
}

func newMultiImageReader(readers map[image.Rectangle]ImageReader) *_MultiImageReader {
	assert(len(readers) > 0)

	p := &_MultiImageReader{
		rect:    image.Rect(0, 0, 1, 1),
		readers: make(map[image.Rectangle]ImageReader),
	}
	for b, r := range readers {
		assertf(!b.Empty(), "b = %v", b)
		assertf(b.Min.X >= 0, "b = %v", b)
		assertf(b.Min.Y >= 0, "b = %v", b)
		assertf(r != nil, "b = %v", b)

		p.readers[b] = r
		p.rect = p.rect.Union(b)
		p.channels = r.Channels()
		p.dataType = r.DataType()
	}
	return p
}

func (p *_MultiImageReader) Close() error {
	var firstErr error
	if len(p.readers) > 0 {
		for _, r := range p.readers {
			if err := r.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		p.readers = nil
	}
	return firstErr
}

func (p *_MultiImageReader) Width() int {
	return p.rect.Dx()
}
func (p *_MultiImageReader) Height() int {
	return p.rect.Dy()
}
func (p *_MultiImageReader) Channels() int {
	return p.channels
}
func (p *_MultiImageReader) DataType() reflect.Kind {
	return p.dataType
}

func (p *_MultiImageReader) HasOverviews() bool {
	if len(p.readers) > 1 {
		return false
	}
	for _, r := range p.readers {
		return r.HasOverviews()
	}
	return false
}
func (p *_MultiImageReader) HasOverviewsFeature() bool {
	if len(p.readers) > 1 {
		return false
	}
	for _, r := range p.readers {
		return r.HasOverviewsFeature()
	}
	return false
}
func (p *_MultiImageReader) BuildOverviews() error {
	if len(p.readers) > 1 {
		return ErrNoOverviewsFeature
	}
	for _, r := range p.readers {
		return r.BuildOverviews()
	}
	return ErrNoOverviewsFeature
}
func (p *_MultiImageReader) BuildOverviewsIfNotExists() error {
	if len(p.readers) > 1 {
		return ErrNoOverviewsFeature
	}
	for _, r := range p.readers {
		return r.BuildOverviewsIfNotExists()
	}
	return ErrNoOverviewsFeature
}
func (p *_MultiImageReader) Read(rect image.Rectangle) (m image.Image, err error) {
	rect = rect.Intersect(p.rect)
	assert(!rect.Empty())

	// only on image, start at (0,0)
	if len(p.readers) == 1 && p.rect.Min == image.Pt(0, 0) {
		for _, r := range p.readers {
			return r.Read(rect)
		}
	}

	// find readers rect
	var rectList []image.Rectangle
	for b, _ := range p.readers {
		if !b.Intersect(rect).Empty() {
			rectList = append(rectList, b)
		}
	}
	if len(rectList) == 0 {
		m = ximage.NewMemPImage(image.Rect(0, 0, rect.Dx(), rect.Dy()), p.channels, p.dataType)
		return m, nil
	}

	// read rect form rectList
	m = ximage.NewMemPImage(image.Rect(0, 0, rect.Dx(), rect.Dy()), p.channels, p.dataType)
	for i := 0; i < len(rectList); i++ {
		r := p.readers[rectList[i]]
		b := rect.Intersect(rectList[i])

		// read sun image
		sub, err := r.Read(b.Sub(rectList[i].Min))
		if err != nil {
			return nil, err
		}

		// copy sub image
		xdraw.Draw(m.(*ximage.MemPImage), b.Sub(rect.Min), sub, image.Pt(0, 0))
	}

	// OK
	return
}
func (p *_MultiImageReader) ReadOverview(idxOverview int, rect image.Rectangle) (m image.Image, err error) {
	assert(idxOverview >= 0)

	rect = rect.Intersect(p.rect)
	assert(!rect.Empty())

	if idxOverview == 0 {
		return p.Read(rect)
	}

	if len(p.readers) != 1 {
		return nil, ErrNoOverviewsFeature
	}
	if !p.HasOverviewsFeature() {
		return nil, ErrNoOverviewsFeature
	}
	if !p.HasOverviews() {
		return nil, ErrNoOverviews
	}

	// only on image, start at (0,0)
	if len(p.readers) == 1 && p.rect.Min == image.Pt(0, 0) {
		for _, r := range p.readers {
			return r.ReadOverview(idxOverview, rect)
		}
	}

	return nil, ErrNoOverviews
}
