// +build darwin freebsd openbsd netbsd dragonfly
// +build !appengine,!js

package formatter

import (
	"io"

	"golang.org/x/sys/unix"
)

const ioctlReadTermios = unix.TIOCGETA

type Termios unix.Termios

func initTerminal(w io.Writer) {
}
