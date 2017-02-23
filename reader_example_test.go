// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image_test

import (
	"fmt"
	"log"

	ximage "github.com/chai2010/image"
)

func ExampleLoad() {
	m, format, err := ximage.Load("./testdata/lena.png")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("format = %s\n", format)
	fmt.Printf("Rect = %v\n", m.Bounds())
	// Output:
	// format = png
	// Rect = (0,0)-(512,512)
}

func ExampleLoadImage() {
	m, format, err := ximage.LoadImage("./testdata/lena.png")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("format = %v\n", format)
	fmt.Printf("MemPMagic = %v\n", m.XMemPMagic)
	fmt.Printf("Rect = %v\n", m.XRect)
	fmt.Printf("Channels = %v\n", m.XChannels)
	fmt.Printf("DataType = %v\n", m.XDataType)
	// Output:
	// format = png
	// MemPMagic = MemP
	// Rect = (0,0)-(512,512)
	// Channels = 4
	// DataType = uint8
}
