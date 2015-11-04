// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package big provides big pyramid image support.
package big

import (
	"errors"
	"image"
	"reflect"
)

var (
	ErrNoOverviews        = errors.New("image/big: no overviews!")
	ErrNoOverviewsFeature = errors.New("image/big: no overviews feature!")
	ErrFormat             = image.ErrFormat
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

type _ImageReaderDriver struct {
	Open func(filename string) (ImageReader, error)
}

type _ImageWriterDriver struct {
	Open   func(filename string) (ImageWriter, error)
	Create func(format, filename string, width, height, channels int, dataType reflect.Kind) (ImageWriter, error)
}

var (
	_ImageReaderDriverList []_ImageReaderDriver
	_ImageWriterDriverList []_ImageWriterDriver
)

func MultiImageReader(readers map[image.Rectangle]ImageReader) ImageReader {
	return &_MultiImageReader{Readers: readers}
}

func MultiOverviewImageReader(readers []map[image.Rectangle]ImageReader) ImageReader {
	return &_MultiOverviewImageReader{Readers: readers}
}

func OpenImageReader(filename string) (r ImageReader, err error) {
	for _, it := range _ImageReaderDriverList {
		if r, err = it.Open(filename); err == nil {
			return
		}
	}
	err = ErrFormat
	return
}

func OpenImageWriter(filename string) (r ImageWriter, err error) {
	for _, it := range _ImageWriterDriverList {
		if r, err = it.Open(filename); err == nil {
			return
		}
	}
	err = ErrFormat
	return
}

func CreateImageWriter(format, filename string, width, height, channels int, dataType reflect.Kind) (w ImageWriter, err error) {
	for _, it := range _ImageWriterDriverList {
		if w, err = it.Create(format, filename, width, height, channels, dataType); err == nil {
			return
		}
	}
	err = ErrFormat
	return
}

func RegisterImageReader(driverName string,
	open func(filename string) (ImageReader, error),
) {
	_ImageReaderDriverList = append(
		_ImageReaderDriverList,
		_ImageReaderDriver{
			Open: open,
		},
	)
}

func RegisterImageWriter(driverName string,
	open func(filename string) (ImageWriter, error),
	create func(format, filename string, width, height, channels int, dataType reflect.Kind) (ImageWriter, error),
) {
	_ImageWriterDriverList = append(
		_ImageWriterDriverList,
		_ImageWriterDriver{
			Open:   open,
			Create: create,
		},
	)
}
