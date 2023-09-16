package atomic

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	procMove    = modkernel.NewProc("MoveFileExW")
)

func moveFileEx(existingFile, newFile *uint16, dwFlags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(
		procMove.Addr(), 3,
		uintptr(unsafe.Pointer(existingFile)),
		uintptr(unsafe.Pointer(newFile)),
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
