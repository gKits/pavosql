//go:build windows

package kvstore

import (
	"errors"
	"os"
	"syscall"
)

var (
	errMmapFileSize = errors.New("mmap: cannot init mmap, file size needs to be multiple of page size")
)

type mmap struct {
	fileSize int
	mmapSize int
	chunks   [][]byte
}

func (mm *mmap) Init(f *os.File) error {
	fStats, err := f.Stat()
	if err != nil {
		return err
	}

	if fStats.Size()%PageSize != 0 {
		return errMmapFileSize
	}

	mmapSize := 64 << 20
	for mmapSize < int(fStats.Size()) {
		mmapSize *= 2
	}

	fileMap, err := syscall.CreateFileMapping(
		syscall.Handle(f.Fd()),
		nil,
		syscall.PAGE_READWRITE,
		0, uint32(mmapSize),
		nil,
	)
	defer syscall.CloseHandle(fileMap)

	addr, err := syscall.MapViewOfFile(fileMap, syscall.FILE_MAP_WRITE, 0, 0, uintptr(mmapSize))
	if err != nil {
		return err
	}
	defer syscall.UnmapViewOfFile(addr)

	// data := (([]*byte)(unsafe.Pointer(addr)))
	mm.mmapSize = mmapSize
	mm.chunks = [][]byte{}
	mm.fileSize = int(fStats.Size())

	return nil
}

func (mm *mmap) Extend(f *os.File, n int) error {
	if mm.mmapSize >= n*PageSize {
		return nil
	}

	// chunk, err := syscall.Mmap(
	// 	int(f.Fd()), int64(mm.mmapSize), mm.mmapSize,
	// 	syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED,
	// )
	// if err != nil {
	// 	return err
	// }

	// mm.mmapSize *= 2
	// mm.chunks = append(mm.chunks, chunk)

	return nil
}

func (mm *mmap) Close() error {
	return nil
}

func (mm *mmap) ExtendFile(f *os.File, n int) error {
	filePages := mm.fileSize / PageSize
	if filePages >= n {
		return nil
	}

	for filePages < n {
		inc := filePages / 8
		if inc < 1 {
			inc = 1
		}
		filePages += inc
	}

	// if err := syscall.LockFile(); err != nil {

	// }

	// fileSize := filePages * pageSize
	// if err := syscall.Fallocate(int(f.Fd()), 0, 0, int64(fileSize)); err != nil {
	// 	return err
	// }

	// mm.fileSize = fileSize
	return nil
}
