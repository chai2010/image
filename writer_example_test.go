// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image_test

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"

	ximage "."
)

func ExampleEncode_default() {
	m, _, err := ximage.Load("./testdata/lena.png")
	if err != nil {
		log.Fatal(err)
	}

	pngOut, err := ximage.Encode(m, nil)
	if err != nil {
		log.Fatal(err)
	}
	_ = pngOut
}

func ExampleEncode_user_defined() {
	m, _, err := ximage.Load("./testdata/lena.png")
	if err != nil {
		log.Fatal(err)
	}

	jpegOut, err := ximage.Encode(m, func(w io.Writer, m image.Image) error {
		return jpeg.Encode(w, m, nil)
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = jpegOut
}

func ExampleSave_png() {
	outfile := "zz_lena.png"
	defer os.Remove(outfile)

	m, format, err := ximage.Load("./testdata/lena.png")
	if err != nil {
		log.Fatal(err)
	}

	err = ximage.Save(outfile, m, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("format = %s\n", format)
	// Output:
	// format = png
}

func ExampleSave_pngAsJpeg() {
	outfile := "zz_lena.png"
	defer os.Remove(outfile)

	m, format, err := ximage.Load("./testdata/lena.png")
	if err != nil {
		log.Fatal(err)
	}

	err = ximage.Save(outfile, m, func(w io.Writer, m image.Image) error {
		return jpeg.Encode(w, m, nil)
	})
	if err != nil {
		log.Fatal(err)
	}

	// outfile is a jpeg
	_, format, err = ximage.Load(outfile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("format = %s\n", format)
	// Output:
	// format = jpeg
}
