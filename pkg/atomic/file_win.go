//go:build windows

package atomic

import (
	"os"
	"syscall"
)

const (
	movefile_replace_existing = 0x1
	movefile_write_through    = 0x8
)

//sys moveFileEx(lpExistingFileName *uint16, lpNewFileName *uint16, dwFlags uint32) (err error) = MoveFileExW

// replaceFile atomically replaces the destination file or directory with the
// source.  It is guaranteed to either replaceFile the target file entirely, or not
// change either file.
func replaceFile(source, destination string) error {
	src, err := syscall.UTF16PtrFromString(source)
	if err != nil {
		return &os.LinkError{Op: "replace", Old: source, New: destination, Err: err}
	}
	dest, err := syscall.UTF16PtrFromString(destination)
	if err != nil {
		return &os.LinkError{Op: "replace", Old: source, New: destination, Err: err}
	}

	// see http://msdn.microsoft.com/en-us/library/windows/desktop/aa365240(v=vs.85).aspx
	if err := moveFileEx(src, dest, movefile_replace_existing|movefile_write_through); err != nil {
		return &os.LinkError{Op: "replace", Old: source, New: destination, Err: err}
	}
	return nil
}
