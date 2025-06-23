package tree

type pager interface {
	ReadPage(off uint64) ([]byte, error)
	Commit() error
	Rollback() error
}

type Tree struct {
	root  uint64
	pager pager
}

func New() *Tree {
	return &Tree{}
}

func (t *Tree) Get(k []byte) ([]byte, error) {
	return nil, nil
}

func (t *Tree) Set(k []byte, v []byte) error {
	return nil
}

func (t *Tree) Delete(k []byte) error {
	return nil
}
