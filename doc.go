// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package image define a Memory format picture.
//
// MemP Image Spec (Native Endian), see https://github.com/chai2010/image.
//	type MemP interface {
//		MemPMagic() string
//		Bounds() image.Rectangle
//		Channels() int
//		DataType() reflect.Kind
//		Pix() []byte // PixSilce type
//
//		// Stride is the Pix stride (in bytes, must align with SizeofKind(p.DataType))
//		// between vertically adjacent pixels.
//		Stride() int
//	}
//
// Please report bugs to chaishushan{AT}gmail.com.
//
// Thanks!
package image
