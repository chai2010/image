// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big

import (
	"image"
	"reflect"
)

var (
	_ ImageReader = (*_MultiImageReader)(nil)
)

type _MultiImageReader struct {
	Readers map[image.Rectangle]ImageReader
}

func (p *_MultiImageReader) Close() error {
	return nil
}

func (p *_MultiImageReader) Width() int {
	return 0
}
func (p *_MultiImageReader) Height() int {
	return 0
}
func (p *_MultiImageReader) Channels() int {
	return 0
}
func (p *_MultiImageReader) DataType() reflect.Kind {
	return 0
}

func (p *_MultiImageReader) HasOverviews() bool {
	return false
}
func (p *_MultiImageReader) HasOverviewsFeature() bool {
	return false
}
func (p *_MultiImageReader) BuildOverviews() error {
	return nil
}
func (p *_MultiImageReader) BuildOverviewsIfNotExists() error {
	return nil
}
func (p *_MultiImageReader) Read(r image.Rectangle) (m image.Image, err error) {
	return
}
func (p *_MultiImageReader) ReadOverview(idxOverview int, r image.Rectangle) (m image.Image, err error) {
	return
}
