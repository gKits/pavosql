//go:build !windows

package backend

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

func (mm *mmap) init(f *os.File) error {
	fStats, err := f.Stat()
	if err != nil {
		return err
	}

	if fStats.Size()%pageSize != 0 {
		return errMmapFileSize
	}

	mmapSize := 64 << 20
	for mmapSize < int(fStats.Size()) {
		mmapSize *= 2
	}

	chunk, err := syscall.Mmap(
		int(f.Fd()), 0, mmapSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED,
	)
	if err != nil {
		return err
	}

	mm.mmapSize = mmapSize
	mm.chunks = [][]byte{chunk}
	mm.fileSize = int(fStats.Size())

	return nil
}

func (mm *mmap) extend(f *os.File, n int) error {
	if mm.mmapSize >= n*pageSize {
		return nil
	}

	chunk, err := syscall.Mmap(
		int(f.Fd()), int64(mm.mmapSize), mm.mmapSize,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED,
	)
	if err != nil {
		return err
	}

	mm.mmapSize *= 2
	mm.chunks = append(mm.chunks, chunk)

	return nil
}

func (mm *mmap) close() error {
	for _, chunk := range mm.chunks {
		if err := syscall.Munmap(chunk); err != nil {
			return err
		}
	}
	return nil
}

func (mm *mmap) extendFile(f *os.File, n int) error {
	filePages := mm.fileSize / pageSize
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

	fileSize := filePages * pageSize
	if err := syscall.Fallocate(int(f.Fd()), 0, 0, int64(fileSize)); err != nil {
		return err
	}

	mm.fileSize = fileSize
	return nil
}
