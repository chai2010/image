// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"reflect"

	ximage "."
)

type BGR struct {
	B, G, R uint16
}

func main() {
	var buf bytes.Buffer
	var data []byte
	var err error

	// Load file data
	if data, err = ioutil.ReadFile("./testdata/lena.jpg"); err != nil {
		log.Println(err)
	}

	// Decode jpeg
	m0, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
	}

	// copy to BGR48 format image
	b := m0.Bounds()
	rgbImage := ximage.NewImage(b, 3, reflect.Uint16)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		var (
			line     []byte = rgbImage.Pix[rgbImage.PixOffset(b.Min.X, y):][:rgbImage.Stride]
			rgbSlice []BGR  = ximage.PixSilce(line).Slice(reflect.TypeOf([]BGR(nil))).([]BGR)
		)

		for x, _ := range rgbSlice {
			r, g, b, _ := m0.At(x, y).RGBA()
			rgbSlice[x] = BGR{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
			}
		}
	}

	// save as jpeg
	if err = jpeg.Encode(&buf, rgbImage, nil); err != nil {
		log.Println(err)
	}
	if err = ioutil.WriteFile("zz_output.jpg", buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	fmt.Println("Done")
}
