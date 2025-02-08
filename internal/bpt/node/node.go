// The node package provides all basic functionalities of a singular node inside a B+Tree. A single
// node can be of two either two types [TypeLeaf] or [TypePointer]. Both node types contain
// key-value pairs with the primary difference beeing that [TypeLeaf] nodes store actual data as
// their values while [TypePointer] store pointers to other nodes that are called their children.
//
// Nodes are stored in a custom byte format that is structured as follows:
//
//	 offset | size in B | description
//	--------+-----------+-------------
//	      0 |         1 | Type identifier
//	      1 |         2 | Number of stored cells N
//	      3 |       2xN | Offset list to each cell
//	  3+2xN |       ... | Cells
//
// The cells that are stored on a [TypeLeaf] node consist of a key and a value. They are also stored
// in the form of a custom byte format with following structure.
//
//	 offset | size in B | description
//	--------+-----------+-------------
//	      0 |         2 | Length of key K
//	      2 |         K | Key
//	    2+K |         2 | Length of value V
//	    4+K |         V | Value
//
// Since a [TypePointer] node only stores uint64 pointers to other nodes their cells structure looks
// slightly different.
//
//	 offset | size in B | description
//	--------+-----------+-------------
//	      0 |         2 | Length of key K
//	      2 |         K | Key
//	    2+K |         8 | Pointer
package node

import (
	"bytes"
	"encoding/binary"
	"slices"
)

type Type uint8

const (
	TypeLeaf    Type = 0x01
	TypePointer Type = 0x10
)

func TypeOf(node []byte) Type {
	t := Type(node[0])
	if t != TypeLeaf && t != TypePointer {
		panic("invalid node type")
	}
	return t
}

func LenOf(node []byte) uint16 {
	return binary.LittleEndian.Uint16(node[1:])
}

func SizeOf(node []byte) uint16 {
	return offset(node, LenOf(node))
}

func Key(node []byte, i uint16) []byte {
	if i >= LenOf(node) {
		panic("index out of bounds")
	}
	off := offset(node, i)
	kLen := binary.LittleEndian.Uint16(node[off:])
	return node[off+2 : off+2+kLen]
}

func Value(node []byte, i uint16) []byte {
	if i >= LenOf(node) {
		panic("index out of bounds")
	}
	off := offset(node, i)
	kLen := binary.LittleEndian.Uint16(node[off:])
	vLen := binary.LittleEndian.Uint16(node[off+2+kLen:])
	return node[off+4+kLen : off+4+kLen+vLen]
}

func Search(node []byte, target []byte) (uint16, bool) {
	n := LenOf(node)
	left, right := uint16(0), n

	for left < right {
		cur := uint16(uint(left+right) >> 1) // #nosec G115 // right shift stops overflow
		if cmp := bytes.Compare(Key(node, cur), target); cmp == -1 {
			left = cur + 1
		} else if cmp == 1 {
			right = cur
		} else {
			return cur, true
		}
	}

	return left, left < n && bytes.Equal(Key(node, left), target)
}

func Insert(node []byte, i uint16, k, v []byte) []byte {
	n := LenOf(node)
	off := offset(node, i)
	cell := makeCell(k, v)

	binary.LittleEndian.PutUint16(node[1:], n+1)

	node = slices.Insert(node, int(off), cell...)

	node = slices.Insert(node, int(3+i*2), 0, 0)
	binary.LittleEndian.PutUint16(node[3+i*2:], off)

	for a := uint16(1); a <= n+1; a++ {
		currentOffset := offset(node, a)
		var shift uint16 = 2
		if a > i {
			shift += uint16(len(cell))
		}
		setOffset(node, a, uint16(currentOffset)+shift)
	}
	return node
}

func Update(node []byte, i uint16, k, v []byte) []byte {
	n := LenOf(node)
	off := offset(node, i)
	cell := makeCell(k, v)

	nextOff := offset(node, i+1)
	node = slices.Replace(node, int(off), int(nextOff), cell...)
	lenDiff := len(cell) - int(nextOff-off)

	for a := i + 1; a <= n; a++ {
		currentOffset := offset(node, a)
		setOffset(node, a, uint16(int(currentOffset)+lenDiff))
	}

	return node
}

func Delete(node []byte, i uint16) []byte {
	n := LenOf(node)
	off := offset(node, i)
	nextOff := offset(node, i+1)

	binary.LittleEndian.PutUint16(node[1:], n-1)
	node = slices.Delete(node, int(off), int(nextOff))

	node = slices.Delete(node, int(3+i*2), int(3+i*2+2))

	for a := uint16(1); a <= n-1; a++ {
		currentOffset := offset(node, a)
		shift := uint16(2)
		if a > i {
			shift += nextOff - off
		}
		setOffset(node, a, uint16(currentOffset-shift))
	}

	return node
}

func Merge(left, right []byte) []byte {
	if TypeOf(left) != TypeOf(right) {
		panic("node types do not match")
	}

	ln := LenOf(left)
	rn := LenOf(right)
	lsz := SizeOf(left)
	rsz := SizeOf(right)

	merged := slices.Insert(left, int(lsz), right[offset(right, 0):rsz]...)
	merged = slices.Insert(merged, int(offset(left, 0)), right[3:3+2*rn]...)
	binary.LittleEndian.PutUint16(merged[1:], ln+rn)

	for i := range LenOf(merged) + 1 {
		curOff := offset(merged, i)
		shift := rn * 2
		if i > ln {
			shift = offset(left, ln) - 3
		}
		setOffset(merged, i, curOff+shift)
	}

	lastOff := offset(left, LenOf(left))

	_ = lastOff

	return merged[:SizeOf(merged)]
}

func Split(node []byte) ([]byte, []byte) {
	center := func(node []byte) uint16 {
		n := LenOf(node)
		target := ((uint16(SizeOf(node)) - offset(node, 0)) >> 1) + offset(node, 0)

		cur := n / 2
		for cur > 0 && cur < n {
			curOff := offset(node, cur)
			prevOff := offset(node, cur-1)
			nextOff := offset(node, cur+1)

			switch {
			case curOff == target:
				return cur
			case nextOff == target:
				return cur + 1
			case prevOff == target:
				return cur - 1
			case curOff < target && nextOff < target:
				cur = cur + 1
				continue
			case curOff > target && prevOff > target:
				cur = cur - 1
			case curOff < target && nextOff > target:
				if target-curOff < nextOff-target {
					return cur
				}
				return cur + 1
			case curOff > target && prevOff < target:
				if curOff-target < target-prevOff {
					return cur
				}
				return cur - 1
			}
		}
		return cur
	}(node)

	n := LenOf(node)
	typ := TypeOf(node)
	off0 := offset(node, 0)
	offC := offset(node, center)

	left := []byte{byte(typ)}
	rigth := []byte{byte(typ)}

	left = binary.LittleEndian.AppendUint16(left, center)
	rigth = binary.LittleEndian.AppendUint16(rigth, n-center)

	left = append(left, node[3:3+2*center]...)
	rigth = append(rigth, node[3+2*center:off0]...)

	left = append(left, node[off0:offC]...)
	rigth = append(rigth, node[offC:]...)

	for i := range LenOf(left) + 1 {
		curOff := offset(left, i)
		setOffset(left, i, curOff-(n-center)*2)
	}
	for i := range LenOf(rigth) + 1 {
		curOff := offset(rigth, i)
		setOffset(rigth, i, curOff-offC+3+2*(n-center))
	}

	return left, rigth
}

func offset(node []byte, i uint16) uint16 {
	if i == 0 {
		return 3 + LenOf(node)*2
	} else if i > LenOf(node) {
		panic("node index out of bounds")
	}
	return binary.LittleEndian.Uint16(node[3+(i-1)*2:])
}

func setOffset(node []byte, i uint16, off uint16) {
	if i > LenOf(node) {
		panic("node index out of bounds")
	} else if i != 0 {
		binary.LittleEndian.PutUint16(node[3+(i-1)*2:], off)
	}
}

func makeCell(k, v []byte) []byte {
	cell := make([]byte, 4+len(k)+len(v))
	binary.LittleEndian.PutUint16(cell[0:], uint16(len(k)))
	copy(cell[2:], k)
	binary.LittleEndian.PutUint16(cell[2+len(k):], uint16(len(v)))
	copy(cell[4+len(k):], v)
	return cell
}
