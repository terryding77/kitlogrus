// Based on ssh/terminal:
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine,!js

package formatter

import (
	"io"

	"golang.org/x/sys/unix"
)

const ioctlReadTermios = unix.TCGETS

type Termios unix.Termios

func initTerminal(w io.Writer) {
}
