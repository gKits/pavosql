package atomic

import (
	"fmt"
	"io"
	"os"
)

type ReadWriterAt interface {
	io.ReaderAt
	io.WriterAt
	Commit() error
	Abort() error
}

type File struct {
	original string
	info     os.FileInfo
	tmp      *os.File
}

func OpenFile(name string) (*File, error) {
	tmp, err := os.CreateTemp("/tmp", "")
	if err != nil {
		return nil, fmt.Errorf("atomic: failed to create temporary file: %w", err)
	}

	og, err := os.Open(name)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("atomic: failed to open target file: %w", err)
	}
	defer og.Close()

	if _, err := io.Copy(tmp, og); err != nil {
		return nil, fmt.Errorf("atomic: failed to copy data into temporary file: %w", err)
	}

	info, err := og.Stat()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("atomic: failed to obtain file info of target file: %w", err)
	}

	// tmp.Chmod

	return &File{
		original: name,
		info:     info,
		tmp:      tmp,
	}, nil
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	return f.tmp.WriteAt(b, off)
}

func (f *File) Commit() error {
	defer os.Remove(f.tmp.Name())
	defer f.tmp.Close()

	if err := f.tmp.Sync(); err != nil {
		return fmt.Errorf("atomic: failed to flush temporary file: %w", err)
	}
	if err := f.tmp.Close(); err != nil {
		return fmt.Errorf("atomic: failed to close temporary file: %w", err)
	}

	if f.info != nil {
		if err := f.tmp.Chmod(f.info.Mode()); err != nil {
			return fmt.Errorf("atomic: failed to set filemode of temporary file: %w", err)
		}
	}

	if err := replaceFile("", f.original); err != nil {
		return fmt.Errorf("atomic: failed to replace target file: %w", err)
	}
	return nil
}

func (f *File) Abort() error {
	defer os.Remove(f.tmp.Name())
	return f.tmp.Close()
}
