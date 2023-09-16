package atomic

import (
	"os"
	"syscall"
)

func ReplaceFile(source, destination string) error {
	src, err := syscall.UTF16PtrFromString(source)
	if err != nil {
		return &os.LinkError{"replace", source, destination, err}
	}
	dest, err := syscall.UTF16PtrFromString(destination)
	if err != nil {
		return &os.LinkError{"replace", source, destination, err}
	}

	if err := moveFile(src, dest, movefile_replace_existing|movefile_write_through); err != nil {
		return &os.LinkError{"replace", source, destination, err}
	}
	return nil
}
