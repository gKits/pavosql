package atomic

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Write(path string, r io.Reader) error {
	dir, file := filepath.Split(path)
	if dir == "" {
		dir = "."
	}

	f, err := os.CreateTemp(dir, file)
	if err != nil {
		return fmt.Errorf("temp file creation failed: %v", err)
	}
	defer func() {
		if err != nil {
			_ = os.Remove(f.Name())
		}
	}()
	defer f.Close()

	fName := f.Name()

	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("write to temp file failed: %v", err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("fsync of temp file failed: %v", err)
	}

	destInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		// og file does not exist
	} else if err != nil {
		return err
	} else {
		sourceInfo, err := os.Stat(file)
		if err != nil {
			return err
		}

		if sourceInfo.Mode() != destInfo.Mode() {
			if err := os.Chmod(file, destInfo.Mode()); err != nil {
				return fmt.Errorf("equalizing filemode of source and temp file failed: %v", err)
			}
		}
	}
	if err := ReplaceFile(fName, path); err != nil {
		return fmt.Errorf("replacing of original file with temp file failed: %v", err)
	}
	return nil
}
