// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package memp define a Memory format picture.
//
// MemP Image Structs (Native Endian):
//	type Image struct {
//		MemPMagic    string // MemP
//		Rect         image.Rectangle
//		Channels     int
//		DataType     reflect.Kind
//		Pix          PixSilce
//
//		// Stride is the Pix stride (in bytes)
//		// between vertically adjacent pixels.
//		Stride int
//	}
//
// Please report bugs to chaishushan{AT}gmail.com.
//
// Thanks!
package memp
