//go: build !windows

package atomic

import (
	"os"
)

func ReplaceFile(source, destination string) error {
	return os.Rename(source, destination)
}
