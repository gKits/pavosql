//go:build windows

package atomic

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procMoveFileExW = modkernel32.NewProc("MoveFileExW")
)

func moveFileEx(lpExistingFileName *uint16, lpNewFileName *uint16, dwFlags uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(
		procMoveFileExW.Addr(),
		3,
		uintptr(unsafe.Pointer(lpExistingFileName)),
		uintptr(unsafe.Pointer(lpNewFileName)),
		uintptr(dwFlags),
	)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
