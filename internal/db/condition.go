package db

type cmp uint8

const (
	cmpEQ cmp = iota
	cmpNEQ
	cmpLT
	cmpLTE
	cmpTween
	cmpLike
	cmpIn
)

type Condition struct {
	col string
	val []byte
	cmp cmp
}

func (cnd *Condition) Check(c Cell) bool {
	return false
}
