// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"reflect"
)

type ImageInfo interface {
	Bounds() image.Rectangle
	DataType() reflect.Kind
	Channels() int
}

type ImageReader interface {
	Read(r image.Rectangle) (m image.Image, err error)
}

type ImageWriter interface {
	Write(r image.Rectangle, m image.Image) error
}

type ImageOverviewInfo interface {
	HasOverviews() bool
	HasOverviewsFeature() bool
}

type ImageOverviewBuilder interface {
	ImageOverviewInfo
	BuildOverviews() error
	BuildOverviewsIfNotExists() error
}

type ImageOverviewReader interface {
	ImageOverviewInfo
	ReadOverview(idxOverview int, r image.Rectangle) (m image.Image, err error)
}
