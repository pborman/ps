//go:build darwin

package ps

import (
	"syscall"
	"unsafe"
)

const (
	_CTL_KERN       = 1
	_KERN_MAXPROC   = 6
	_KERN_PROC      = 14
	_KERN_PROCARGS2 = 49
	_KERN_PROC_ALL  = 0
	_KERN_PROC_PID  = 1
)

func sysctl(mib []int32) ([]byte, error) {
	var data []byte
	var size int
	var err error

	// Keep looping until a call has sufficient data to get the entire
	// table.  We start with no data which will return the full table
	// size.  It is possible the table size will increase between
	// calls so keep looping until we get them all.
	for {
		size, data, err = sysctl1(mib, data)
		if err == nil && size > len(data) {
			data = make([]byte, size)
			continue
		}
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func sysctl1(mib []int32, data []byte) (int, []byte, error) {
	size := uintptr(len(data))
	var buffer uintptr
	if size > 0 {
		buffer = uintptr(unsafe.Pointer(&data[0]))
	}
	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(len(mib)),
		buffer,
		uintptr(unsafe.Pointer(&size)),
		0,
		0)
	if errno != 0 {
		return 0, nil, errno
	}
	if int(size) < len(data) {
		data = data[:size]
	}
	return int(size), data, nil
}
