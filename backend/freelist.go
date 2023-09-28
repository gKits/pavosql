package backend

import "errors"

type getFuncFN func(uint64) (freelistNode, error)
type pullFuncFN func(uint64) (freelistNode, error)

type freelist struct {
	root  uint64     // pointer to root node
	get   getFuncFN  // callback to get node
	pull  pullFuncFN // callback to get node and free the pointer
	alloc allocFunc  // callback to allocate node and get the new pointer
	free  freeFunc   // callback to free node space at given pointer
}

var (
	errFLPopEmpty = errors.New("freelist: cannot Pop from empty freelist")
)

func (fl *freelist) Pop() (uint64, error) {
	if fl.root == 0 {
		return 0, errFLPopEmpty
	}

	root, err := fl.pull(fl.root)
	if err != nil {
		return 0, err
	}

	root, res, err := fl.freelistPop(root)
	if err != nil {
		return 0, err
	}

	fl.root, err = fl.alloc(root)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (fl *freelist) Push(val uint64) error {
	var root freelistNode
	var err error

	if fl.root == 0 {
		root = freelistNode{}
	} else {
		root, err = fl.pull(fl.root)
		if err != nil {
			return err
		}
	}

	root, err = fl.freelistPush(root, val)
	if err != nil {
		return err
	}

	fl.root, err = fl.alloc(root)
	if err != nil {
		return err
	}

	return nil
}

func (fl *freelist) freelistPop(fn freelistNode) (freelistNode, uint64, error) {
	if fn.next == 0 {
		ptr, popped := fn.Pop()
		return popped, ptr, nil
	}

	next, err := fl.pull(fn.next)
	if err != nil {
		return freelistNode{}, 0, err
	}

	next, res, err := fl.freelistPop(next)
	if err != nil {
		return freelistNode{}, 0, err
	}

	// free next when it does not contain any pointers
	if next.Total() == 0 {
		if err := fl.free(fn.next); err != nil {
			return freelistNode{}, 0, err
		}
		fn.next = 0
	} else {
		fn.next, err = fl.alloc(next)
		if err != nil {
			return freelistNode{}, 0, err
		}
	}

	return fn, res, nil
}

func (fl *freelist) freelistPush(fn freelistNode, val uint64) (freelistNode, error) {
	if fn.next == 0 {
		// create sub node when added pointer would exceed page size
		if fn.Size()+8 >= PageSize {
			sub := freelistNode{0, []uint64{val}}

			ptr, err := fl.alloc(sub)
			if err != nil {
				return freelistNode{}, err
			}

			fn.next = ptr
			return fn, nil
		}

		pushed := fn.Push(val)
		return pushed, nil
	}

	next, err := fl.pull(fn.next)
	if err != nil {
		return freelistNode{}, err
	}

	next, err = fl.freelistPush(next, val)
	if err != nil {
		return freelistNode{}, err
	}

	fn.next, err = fl.alloc(next)
	if err != nil {
		return freelistNode{}, err
	}

	return fn, nil
}
