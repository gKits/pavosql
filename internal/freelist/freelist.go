package freelist

import ()

type getFunc func(uint64) ([]byte, error)
type pullFunc func(uint64) ([]byte, error)
type allocFunc func([]byte) (uint64, error)
type freeFunc func(uint64) error

type FreelistData struct {
	head     uint64
	ptrs     []uint64
	nHead    int
	nDiscard int
}

type Freelist struct {
	FreelistData
	version uint64
	minRead uint64
	pgSize  int
	get     func(uint64) (freelistNode, error)
	pull    func(uint64) (freelistNode, error)
	alloc   func(freelistNode) (uint64, error)
	free    func(uint64) error
}

func New(
	head, version uint64, pgSize int,
	get getFunc, pull pullFunc, alloc allocFunc, free freeFunc,
) Freelist {
	fl := Freelist{
		version: version,
		pgSize:  pgSize,
		get: func(ptr uint64) (freelistNode, error) {
			d, err := get(ptr)
			if err != nil {
				return freelistNode{}, err
			}
			return decodeFreelistNode(d), nil
		},
		pull: func(ptr uint64) (freelistNode, error) {
			d, err := pull(ptr)
			if err != nil {
				return freelistNode{}, err
			}
			return decodeFreelistNode(d), nil
		},
		alloc: func(fn freelistNode) (uint64, error) {
			return alloc(fn.Encode())
		},
		free: free,
	}

	fl.head = head

	return fl
}

func (fl *Freelist) Nq(ptr uint64) error {
	var head freelistNode
	var err error

	if fl.head == 0 {
		head = freelistNode{next: fl.head}
		head = head.Nq(ptr)
		goto alloc
	}

	head, err = fl.get(fl.head)
	if err != nil {
		return err
	}

	if head.Size()+8 <= fl.pgSize {
		head = head.Nq(ptr)
		if err := fl.free(fl.head); err != nil {
			return err
		}
		goto alloc
	}
	head = freelistNode{next: fl.head}
	head = head.Nq(ptr)

alloc:
	fl.head, err = fl.alloc(head)
	return err
}

func (fl *Freelist) Dq() (uint64, error) {
	if fl.head == 0 {
		return 0, nil
	}

	var node freelistNode
	var err error

	for next := fl.head; next != 0; next = node.next {
		node, err = fl.get(next)
		if err != nil {
			return 0, err
		}
	}

	var ptr uint64
	ptr, node = node.Dq()

	return ptr, nil
}
