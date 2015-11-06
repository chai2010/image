// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"errors"
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
	if len(readers) == 0 {
		return &_MultiOverviewImageReader{}
	}

	p := &_MultiOverviewImageReader{}
	for i := 0; i < len(readers); i++ {
		r := newMultiImageReader(readers[i])
		p.readers = append(p.readers, r)
	}
	return p
}

func (p *_MultiOverviewImageReader) Close() error {
	if len(p.readers) == 0 {
		return nil
	}
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
	if len(p.readers) == 0 {
		return 0
	}
	return p.readers[0].Width()
}
func (p *_MultiOverviewImageReader) Height() int {
	if len(p.readers) == 0 {
		return 0
	}
	return p.readers[0].Height()
}
func (p *_MultiOverviewImageReader) Channels() int {
	if len(p.readers) == 0 {
		return 0
	}
	return p.readers[0].Channels()
}
func (p *_MultiOverviewImageReader) DataType() reflect.Kind {
	if len(p.readers) == 0 {
		return reflect.Invalid
	}
	return p.readers[0].DataType()
}

func (p *_MultiOverviewImageReader) HasOverviews() bool {
	if len(p.readers) == 0 {
		return false
	}
	return len(p.readers) > 1
}
func (p *_MultiOverviewImageReader) HasOverviewsFeature() bool {
	if len(p.readers) == 0 {
		return false
	}
	return len(p.readers) > 1
}
func (p *_MultiOverviewImageReader) BuildOverviews() error {
	if len(p.readers) == 0 {
		return errors.New("image/big: _MultiOverviewImageReader.BuildOverviews, no reader!")
	}
	return p.readers[len(p.readers)-1].BuildOverviews()
}
func (p *_MultiOverviewImageReader) BuildOverviewsIfNotExists() error {
	if len(p.readers) == 0 {
		return errors.New("image/big: _MultiOverviewImageReader.BuildOverviewsIfNotExists, no reader!")
	}
	return p.readers[len(p.readers)-1].BuildOverviewsIfNotExists()
}
func (p *_MultiOverviewImageReader) Read(rect image.Rectangle) (m image.Image, err error) {
	if len(p.readers) == 0 {
		return nil, errors.New("image/big: _MultiOverviewImageReader.Read, no reader!")
	}
	if rect.Empty() {
		return nil, errors.New("image/big: _MultiOverviewImageReader.Read, empty rect!")
	}
	return p.readers[0].Read(rect)
}
func (p *_MultiOverviewImageReader) ReadOverview(idxOverview int, rect image.Rectangle) (m image.Image, err error) {
	if len(p.readers) == 0 {
		return nil, errors.New("image/big: _MultiOverviewImageReader.ReadOverview, no reader!")
	}
	if idxOverview < 0 {
		return nil, errors.New("image/big: _MultiOverviewImageReader.ReadOverview, invalid idxOverview!")
	}
	if rect.Empty() {
		return nil, errors.New("image/big: _MultiOverviewImageReader.ReadOverview, empty rect!")
	}

	if idxOverview < len(p.readers) {
		return p.readers[idxOverview].Read(rect)
	} else {
		return p.readers[idxOverview].ReadOverview(
			idxOverview-len(p.readers)-1,
			rect,
		)
	}
}
