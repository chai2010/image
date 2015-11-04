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
	Readers []map[image.Rectangle]ImageReader
}

func (p *_MultiOverviewImageReader) Close() error {
	return nil
}

func (p *_MultiOverviewImageReader) Width() int {
	return 0
}
func (p *_MultiOverviewImageReader) Height() int {
	return 0
}
func (p *_MultiOverviewImageReader) Channels() int {
	return 0
}
func (p *_MultiOverviewImageReader) DataType() reflect.Kind {
	return 0
}

func (p *_MultiOverviewImageReader) HasOverviews() bool {
	return false
}
func (p *_MultiOverviewImageReader) HasOverviewsFeature() bool {
	return false
}
func (p *_MultiOverviewImageReader) BuildOverviews() error {
	return nil
}
func (p *_MultiOverviewImageReader) BuildOverviewsIfNotExists() error {
	return nil
}
func (p *_MultiOverviewImageReader) Read(r image.Rectangle) (m image.Image, err error) {
	return
}
func (p *_MultiOverviewImageReader) ReadOverview(idxOverview int, r image.Rectangle) (m image.Image, err error) {
	return
}
