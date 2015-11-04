// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package big provides big pyramid image support.
package big

import (
	"image"
	"reflect"
)

type Image interface {
	Width() int
	Height() int
	Channels() int
	DataType() reflect.Kind
	Close() error
}

type ImageReader interface {
	Image
	HasOverviews() bool
	HasOverviewsFeature() bool
	BuildOverviews() error
	BuildOverviewsIfNotExists() error
	Read(r image.Rectangle) (m image.Image, err error)
	ReadOverview(idxOverview int, r image.Rectangle) (m image.Image, err error)
}

type ImageWriter interface {
	Image
	Write(r image.Rectangle, m image.Image) error
}

func MultiImageReader(readers map[image.Rectangle]ImageReader) ImageReader {
	panic("TODO")
}

func MultiOverviewImageReader(readers []map[image.Rectangle]ImageReader) ImageReader {
	panic("TODO")
}

func OpenImageReader(filename string) (r ImageReader, err error) {
	panic("TODO")
}

func OpenImageWriter(filename string) (r ImageWriter, err error) {
	panic("TODO")
}

func CreateImageWriter(format, filename string, width, height, channels int, dataType reflect.Kind) (r ImageWriter, err error) {
	panic("TODO")
}

func RegisterImageReader(driverName string,
	open func(filename string) (ImageReader, error),
) {
	panic("TODO")
}

func RegisterImageWriter(driverName string,
	open func(filename string) (ImageWriter, error),
	create func(format, filename string, width, height, channels int, dataType reflect.Kind) (r ImageWriter, err error),
) {
	panic("TODO")
}
