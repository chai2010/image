// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"image"
	"reflect"
)

var (
	_ ImageReader = (*_MultiOverviewImageReader)(nil)
)

type _MultiOverviewImageReader struct {
	readers []*_MultiImageReader
}

func newMultiOverviewImageReader(readers []map[image.Rectangle]ImageReader) *_MultiOverviewImageReader {
	assert(len(readers) > 0)

	p := &_MultiOverviewImageReader{}
	for i := 0; i < len(readers); i++ {
		r := newMultiImageReader(readers[i])
		p.readers = append(p.readers, r)
	}
	return p
}

func (p *_MultiOverviewImageReader) Close() error {
	var firstErr error
	for i := 0; i < len(p.readers); i++ {
		if err := p.readers[i].Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	p.readers = nil
	return firstErr
}

func (p *_MultiOverviewImageReader) Width() int {
	return p.readers[0].Width()
}
func (p *_MultiOverviewImageReader) Height() int {
	return p.readers[0].Height()
}
func (p *_MultiOverviewImageReader) Channels() int {
	return p.readers[0].Channels()
}
func (p *_MultiOverviewImageReader) DataType() reflect.Kind {
	return p.readers[0].DataType()
}

func (p *_MultiOverviewImageReader) HasOverviews() bool {
	return len(p.readers) > 1
}
func (p *_MultiOverviewImageReader) HasOverviewsFeature() bool {
	return len(p.readers) > 1
}
func (p *_MultiOverviewImageReader) BuildOverviews() error {
	return p.readers[len(p.readers)-1].BuildOverviews()
}
func (p *_MultiOverviewImageReader) BuildOverviewsIfNotExists() error {
	return p.readers[len(p.readers)-1].BuildOverviewsIfNotExists()
}
func (p *_MultiOverviewImageReader) Read(r image.Rectangle) (m image.Image, err error) {
	return p.readers[0].Read(r)
}
func (p *_MultiOverviewImageReader) ReadOverview(idxOverview int, r image.Rectangle) (m image.Image, err error) {
	assert(idxOverview >= 0)

	if idxOverview < len(p.readers) {
		return p.readers[idxOverview].Read(r)
	} else {
		return p.readers[idxOverview].ReadOverview(
			idxOverview-len(p.readers)-1,
			r,
		)
	}
	return
}
