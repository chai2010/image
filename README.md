# Image: Memory Picture

[![Build Status](https://travis-ci.org/chai2010/image.svg)](https://travis-ci.org/chai2010/image)
[![GoDoc](https://godoc.org/github.com/chai2010/image?status.svg)](https://godoc.org/github.com/chai2010/image)

MemP Image Spec (Native Endian):

```Go
// MemP Image Spec (Native Endian), see https://github.com/chai2010/image.
type MemP interface {
	MemPMagic() string
	Bounds() image.Rectangle
	Channels() int
	DataType() reflect.Kind
	Pix() []byte // PixSlice type

	// Stride is the Pix stride (in bytes, must align with SizeofKind(p.DataType))
	// between vertically adjacent pixels.
	Stride() int
}
```

PkgDoc: [http://godoc.org/github.com/chai2010/image](http://godoc.org/github.com/chai2010/image)

Install
=======

1. `go get github.com/chai2010/image`
2. `go run hello.go`


Example
=======

This is a simple example:

```Go
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
```

BUGS
====

Report bugs to <chaishushan@gmail.com>.

Thanks!
