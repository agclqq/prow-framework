package disk

import "syscall"

func Block() int {
	bs := syscall.Getpagesize()
	return bs
}
