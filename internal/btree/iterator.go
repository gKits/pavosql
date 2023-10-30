package btree

type Iterator struct {
	bt    *BTree
	path  []int
	nodes []node
}

func (iter *Iterator) Next() bool {
	depth := len(iter.path) - 1

	if iter.path[depth]+1 >= iter.nodes[depth].Total() {
		iter.path = iter.path[:depth]
		iter.nodes = iter.nodes[:depth]
		return iter.Next()
	}

	switch iter.nodes[depth].Type() {
	case btreePointer:
		ptr := iter.nodes[depth].(pointerNode)

		next, err := iter.bt.get(ptr.ptrs[iter.path[depth]+1])
		if err != nil {
			return false
		}

		iter.nodes = append(iter.nodes, next)
		iter.path = append(iter.path, -1)
		return iter.Next()

	case btreeLeaf:
		iter.path[len(iter.path)-1]++
		return true
	}

	return false
}

func (iter *Iterator) Read() (k []byte, v []byte) {
	leaf := iter.nodes[len(iter.nodes)-1].(leafNode)
	return leaf.keys[iter.path[len(iter.path)-1]], leaf.vals[iter.path[len(iter.path)-1]]
}
