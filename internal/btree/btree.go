package btree

import (
	"io"

	"github.com/gKits/PavoSQL/internal/btree/node"
)

type BTree struct {
	root     uint64
	nodeSize int
	reader   io.ReaderAt
	writer   io.WriterAt
}

func (bt *BTree) Get(key []byte) ([]byte, error) {
	cur, _ := bt.getNode(bt.root)

	for {
		i, exists := cur.Search(key)

		switch cur.Type() {
		case node.TypeLeaf:
			if !exists {
				return nil, nil
			}
			leaf, ok := cur.(*node.Leaf)
			if !ok {
				return nil, nil
			}
			val, err := leaf.Val(i)
			if err != nil {
				return nil, err
			}
			return val, nil

		case node.TypePointer:
			pointer, ok := cur.(*node.Pointer)
			if !ok {
				return nil, nil
			}
			ptr, err := pointer.Ptr(i)
			if err != nil {
				return nil, err
			}
			cur, err = bt.getNode(ptr)
			if err != nil {
				return nil, err
			}
			continue

		default:
			return nil, nil
		}
	}
}

func (bt *BTree) Insert(key, val []byte) error {
	cur, _ := bt.getNode(bt.root)

	for {
		i, exists := cur.Search(key)
		if exists {
			return nil
		}

		switch cur.Type() {
		case node.TypeLeaf:
			leaf, ok := cur.(*node.Leaf)
			if !ok {
				return nil
			}
			if err := leaf.Insert(i, key, val); err != nil {
				return err
			}

		case node.TypePointer:
			pointer, ok := cur.(*node.Pointer)
			if !ok {
				return nil
			}
			ptr, err := pointer.Ptr(i)
			if err != nil {
				return err
			}
			cur, err = bt.getNode(ptr)
			if err != nil {
				return err
			}
			continue

		default:
			return nil
		}
		break
	}

	if cur.Size() <= bt.nodeSize {
		return nil
	}

	// TODO: split node and update path

	return nil
}

func (bt *BTree) getNode(ptr uint64) (noder, error) {
	b := make([]byte, bt.nodeSize)
	if n, err := bt.reader.ReadAt(b, int64(ptr)); err != nil {
		return nil, err
	} else if n != bt.nodeSize {
		return nil, nil
	}

	switch node.TypeOf(b) {
	case node.TypePointer:
		return node.NewPointer(b)
	case node.TypeLeaf:
		return node.NewLeaf(b)
	default:
		return nil, nil
	}
}

func (bt *BTree) writeNode(n noder, ptr uint64) error {
	b := make([]byte, bt.nodeSize)
	nB, err := n.Bytes()
	if err != nil {
		return err
	}
	copy(b, nB)
	if n, err := bt.writer.WriteAt(b, int64(ptr)); err != nil {
		return err
	} else if n != bt.nodeSize {
		return nil
	}

	return nil
}

func (bt *BTree) freeNode(ptr uint64) error {
	if n, err := bt.writer.WriteAt(make([]byte, bt.nodeSize), int64(ptr)); err != nil {
		return err
	} else if n != bt.nodeSize {
		return nil
	}
	return nil
}
