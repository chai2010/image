// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Make pyramid for image file if not exists.
//
//	Usage: mkpyd infile outfile [AB|NN]
//	       mkpyd -h
//
//	Example:
//	  mkpyd infile filename.pyd.rdb
//	  mkpyd infile filename.pyd.rdb AB
//	  mkpyd infile filename.pyd.rdb NN
//
// AB(ApproxBiLinear) is a mixture of the nearest neighbor and bi-linear
// interpolators. It is fast, but usually gives medium quality results.
//
// NN(NearestNeighbor) is the nearest neighbor interpolator. It is very fast,
// but usually gives very low quality results. When scaling up, the result
// will look 'blocky'.
//
// Default is AB(ApproxBiLinear).
//
//	Report bugs to <chaishushan{AT}gmail.com>.
//
package main

import (
	"fmt"
	"os"

	_ "github.com/chai2010/cache"
	_ "github.com/chai2010/gdal"
	_ "github.com/chai2010/image/draw"
	_ "github.com/chai2010/image/pyramid"
	_ "github.com/chai2010/rawp"
	_ "github.com/chai2010/rdb"
	_ "github.com/chai2010/webp"
)

const usage = `
Usage: mkpyd filename outfile [AB|NN]
       mkpyd -h

Example:
  mkpyd filename filename.pyd.rdb
  mkpyd filename filename.pyd.rdb AB
  mkpyd filename filename.pyd.rdb NN

AB(ApproxBiLinear) is a mixture of the nearest neighbor and bi-linear
interpolators. It is fast, but usually gives medium quality results.

NN(NearestNeighbor) is the nearest neighbor interpolator. It is very fast,
but usually gives very low quality results. When scaling up, the result
will look 'blocky'.

Default is AB(ApproxBiLinear).

Report bugs to <chaishushan{AT}gmail.com>.
`

func main() {
	if len(os.Args) < 3 || os.Args[1] == "-h" {
		fmt.Fprintln(os.Stderr, usage[1:len(usage)-1])
		os.Exit(0)
	}

	infile, outfile, resampleTypeName := os.Args[1], os.Args[2], "AB"
	if len(os.Args) > 3 {
		resampleTypeName = os.Args[3]
	}

	_ = infile
	_ = outfile
	_ = resampleTypeName

	panic("TODO")
}
