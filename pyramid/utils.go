// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pyramid

func minInt(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}
