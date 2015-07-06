// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package big

import (
	"fmt"
	"io"
	"runtime"
)

func (defFS) Lock(name string) (io.Closer, error) {
	return nil, fmt.Errorf("big: file locking is not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
