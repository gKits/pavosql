package node

type NodeType uint16

const (
	PNTR_NODE NodeType = iota
	LEAF_NODE
	FLST_NODE
)

type Node interface {
	Type() NodeType
	NKeys() int
	Decode([]byte) error
	Encode() []byte
	Size() int
	Key(int) ([]byte, error)
	Search([]byte) (int, bool)
	Merge(Node) (Node, error)
	Split() (Node, Node)
}
