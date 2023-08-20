package btree

import "fmt"

var (
	ErrWrongCellSize      = &bTreeError{"cell size does not match pages predefined cell size"}
	ErrCellIndexOutBounds = &bTreeError{"cell index exceeds the page's number of cells"}
)

type bTreeError struct {
	Msg string
}

func (bte *bTreeError) Error() string {
	return fmt.Sprintf("btree: %s", bte.Msg)
}
