package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeOf(t *testing.T) {
	cases := map[string]struct {
		node        []byte
		want        Type
		assertPanic assert.PanicAssertionFunc
	}{
		"get leaf node type": {
			node:        []byte{byte(TypeLeaf)},
			want:        TypeLeaf,
			assertPanic: assert.NotPanics,
		},
		"get pointer node type": {
			node:        []byte{byte(TypePointer)},
			want:        TypePointer,
			assertPanic: assert.NotPanics,
		},
		"panic on invalid node type 0x00": {
			node:        []byte{0x00},
			assertPanic: assert.Panics,
		},
		"panic on invalid node type 0x11": {
			node:        []byte{0x11},
			assertPanic: assert.Panics,
		},
		"panic on invalid node type 0xff": {
			node:        []byte{0xff},
			assertPanic: assert.Panics,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c.assertPanic(t, func() {
				got := TypeOf(c.node)
				assert.Equal(t, c.want, got)
			})
		})
	}
}

func TestKey(t *testing.T) {
	cases := map[string]struct {
		node        []byte
		input       uint16
		want        []byte
		assertPanic assert.PanicAssertionFunc
	}{
		"0th key": {
			node:        fixtureNodeEvenNumberOfKeys(),
			input:       0,
			want:        []byte("aaa"),
			assertPanic: assert.NotPanics,
		},
		"1st key": {
			node:        fixtureNodeEvenNumberOfKeys(),
			input:       1,
			want:        []byte("bbb"),
			assertPanic: assert.NotPanics,
		},
		"9th key": {
			node:        fixtureNodeEvenNumberOfKeys(),
			input:       9,
			want:        []byte("jjj"),
			assertPanic: assert.NotPanics,
		},
		"10th key": {
			node:        fixtureNodeEvenNumberOfKeys(),
			input:       10,
			assertPanic: assert.Panics,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c.assertPanic(t, func() {
				got := Key(c.node, c.input)
				assert.Equal(t, c.want, got)
			})
		})
	}
}

func TestValue(t *testing.T) {
	cases := map[string]struct {
		node        []byte
		input       uint16
		want        []byte
		assertPanic assert.PanicAssertionFunc
	}{
		"0th value": {
			node:        fixtureNodeEvenNumberOfKeys(),
			input:       0,
			want:        []byte("aaa"),
			assertPanic: assert.NotPanics,
		},
		"1st value": {
			node:        fixtureNodeEvenNumberOfKeys(),
			input:       1,
			want:        []byte("bbb"),
			assertPanic: assert.NotPanics,
		},
		"9th value": {
			node:        fixtureNodeEvenNumberOfKeys(),
			input:       9,
			want:        []byte("jjj"),
			assertPanic: assert.NotPanics,
		},
		"10th value": {
			node:        fixtureNodeEvenNumberOfKeys(),
			input:       10,
			assertPanic: assert.Panics,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c.assertPanic(t, func() {
				got := Value(c.node, c.input)
				assert.Equal(t, c.want, got)
			})
		})
	}
}

func TestSearch(t *testing.T) {
	cases := map[string]struct {
		node, input  []byte
		want         uint16
		assertExists assert.BoolAssertionFunc
	}{
		"search existing target in odd length node": {
			node:         fixtureNodeOddNumberOfKeys(),
			input:        []byte("ddd"),
			want:         3,
			assertExists: assert.True,
		},
		"search non existing target in odd length node": {
			node:         fixtureNodeOddNumberOfKeys(),
			input:        []byte("iij"),
			want:         9,
			assertExists: assert.False,
		},
		"search existing target in even length node": {
			node:         fixtureNodeEvenNumberOfKeys(),
			input:        []byte("hhh"),
			want:         7,
			assertExists: assert.True,
		},
		"search non existing target in even length node": {
			node:         fixtureNodeEvenNumberOfKeys(),
			input:        []byte("bbc"),
			want:         2,
			assertExists: assert.False,
		},
		"search before first key in odd length node": {
			node:         fixtureNodeOddNumberOfKeys(),
			input:        []byte("aa"),
			want:         0,
			assertExists: assert.False,
		},
		"search after last key in odd length node": {
			node:         fixtureNodeOddNumberOfKeys(),
			input:        []byte("zzzzz"),
			want:         11,
			assertExists: assert.False,
		},
		"search before first key in even length node": {
			node:         fixtureNodeEvenNumberOfKeys(),
			input:        []byte("aa"),
			want:         0,
			assertExists: assert.False,
		},
		"search after last key in even length node": {
			node:         fixtureNodeEvenNumberOfKeys(),
			input:        []byte("zzzzz"),
			want:         10,
			assertExists: assert.False,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, exists := Search(c.node, c.input)
			assert.Equal(t, c.want, got)
			c.assertExists(t, exists)
		})
	}
}

func TestInsert(t *testing.T) {
	cases := map[string]struct {
		node []byte
		i    uint16
		k, v []byte
		want []byte

		wantLen    uint16
		wantKAfter []byte
	}{
		"insert new cell": {
			i: 1, k: []byte("new"), v: []byte("new"),
			node: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
			want: []byte{
				1, 4, 0, 17, 0, 27, 0, 33, 0, 39, 0,
				1, 0, 'a', 1, 0, 'a',
				3, 0, 'n', 'e', 'w', 3, 0, 'n', 'e', 'w',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
		},
		"insert new cell at the start": {
			i: 0, k: []byte("key"), v: []byte("value"),
			node: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
			want: []byte{
				1, 4, 0, 23, 0, 29, 0, 35, 0, 41, 0,
				3, 0, 'k', 'e', 'y', 5, 0, 'v', 'a', 'l', 'u', 'e',
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
		},
		"insert new cell at the end": {
			i: 2, k: []byte("at"), v: []byte("the end"),
			node: []byte{
				1, 2, 0, 13, 0, 19, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
			},
			want: []byte{
				1, 3, 0, 15, 0, 21, 0, 34, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				2, 0, 'a', 't', 7, 0, 't', 'h', 'e', ' ', 'e', 'n', 'd',
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := Insert(c.node, c.i, c.k, c.v)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := map[string]struct {
		node []byte
		i    uint16
		k, v []byte
		want []byte
	}{
		"update a single cell node": {
			i: 0, k: []byte("a"), v: []byte("this is a new value"),
			node: []byte{
				1, 1, 0, 11, 0,
				1, 0, 'a', 1, 0, 'a',
			},
			want: []byte{
				1, 1, 0, 29, 0,
				1, 0, 'a', 19, 0, 't', 'h', 'i', 's', ' ', 'i', 's', ' ', 'a', ' ', 'n', 'e', 'w', ' ', 'v', 'a', 'l', 'u', 'e',
			},
		},
		"update first cell of a multi cell node": {
			i: 0, k: []byte("a"), v: []byte("new"),
			node: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
			want: []byte{
				1, 3, 0, 17, 0, 23, 0, 29, 0,
				1, 0, 'a', 3, 0, 'n', 'e', 'w',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
		},
		"update last cell of a multi cell node": {
			i: 2, k: []byte("ccc"), v: []byte("new"),
			node: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
			want: []byte{
				1, 3, 0, 15, 0, 21, 0, 31, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				3, 0, 'c', 'c', 'c', 3, 0, 'n', 'e', 'w',
			},
		},
		"update cell with shorter cell": {
			i: 2, k: []byte("c"), v: []byte("c"),
			node: []byte{
				1, 3, 0, 15, 0, 21, 0, 29, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				2, 0, 'c', 'c', 2, 0, 'c', 'c',
			},
			want: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := Update(c.node, c.i, c.k, c.v)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := map[string]struct {
		node []byte
		i    uint16
		want []byte
	}{
		"delete last cell of single cell node": {
			i: 0,
			node: []byte{
				1, 1, 0, 11, 0,
				1, 0, 'a', 1, 0, 'a',
			},
			want: []byte{
				1, 0, 0,
			},
		},
		"delete first cell of multi cell node": {
			i: 0,
			node: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
			want: []byte{
				1, 2, 0, 13, 0, 19, 0,
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
		},
		"delete last cell of multi cell node": {
			i: 2,
			node: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
			want: []byte{
				1, 2, 0, 13, 0, 19, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
			},
		},
		"delete cell of multi cell node": {
			i: 1,
			node: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c',
			},
			want: []byte{
				1, 2, 0, 13, 0, 19, 0,
				1, 0, 'a', 1, 0, 'a',
				1, 0, 'c', 1, 0, 'c',
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := Delete(c.node, c.i)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestMerge(t *testing.T) {
	cases := map[string]struct {
		left, right []byte
		want        []byte
	}{
		"merge two equally sized nodes": {
			left: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b', 1, 0, 'c', 1, 0, 'c',
			},
			right: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'd', 1, 0, 'd', 1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
			},
			want: []byte{
				1, 6, 0, 21, 0, 27, 0, 33, 0, 39, 0, 45, 0, 51, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
			},
		},
		"merge two differently sized nodes": {
			left: []byte{
				1, 4, 0, 17, 0, 23, 0, 29, 0, 35, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b', 1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
			},
			right: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f', 1, 0, 'g', 1, 0, 'g',
			},
			want: []byte{
				1, 7, 0, 23, 0, 29, 0, 35, 0, 41, 0, 47, 0, 53, 0, 59, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
				1, 0, 'g', 1, 0, 'g',
			},
		},
		"merge two differently sized nodes with trailing zeros": {
			left: []byte{
				1, 4, 0, 17, 0, 23, 0, 29, 0, 35, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b', 1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
				0, 0, 0, 0, 0,
			},
			right: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f', 1, 0, 'g', 1, 0, 'g',
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			want: []byte{
				1, 7, 0, 23, 0, 29, 0, 35, 0, 41, 0, 47, 0, 53, 0, 59, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
				1, 0, 'g', 1, 0, 'g',
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := Merge(c.left, c.right)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestSplit(t *testing.T) {
	cases := map[string]struct {
		node         []byte
		wantL, wantR []byte
	}{
		"split symmetrical node with even number of cells": {
			node: []byte{
				1, 6, 0, 21, 0, 27, 0, 33, 0, 39, 0, 45, 0, 51, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
			},
			wantL: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b', 1, 0, 'c', 1, 0, 'c',
			},
			wantR: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'd', 1, 0, 'd', 1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
			},
		},
		"split symmetrical node with odd number of cells": {
			node: []byte{
				1, 7, 0, 23, 0, 29, 0, 35, 0, 41, 0, 47, 0, 53, 0, 59, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
				1, 0, 'g', 1, 0, 'g',
			},
			wantL: []byte{
				1, 4, 0, 17, 0, 23, 0, 29, 0, 35, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b', 1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
			},
			wantR: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f', 1, 0, 'g', 1, 0, 'g',
			},
		},
		"split asymmetrical node with even number of cells": {
			node: []byte{
				1, 6, 0, 22, 0, 28, 0, 36, 0, 42, 0, 48, 0, 54, 0,
				2, 0, 'a', 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b',
				2, 0, 'c', 'c', 2, 0, 'c', 'c', 1, 0, 'd', 1, 0, 'd',
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
			},
			wantL: []byte{
				1, 3, 0, 16, 0, 22, 0, 30, 0,
				2, 0, 'a', 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b', 2, 0, 'c', 'c', 2, 0, 'c', 'c',
			},
			wantR: []byte{
				1, 3, 0, 15, 0, 21, 0, 27, 0,
				1, 0, 'd', 1, 0, 'd', 1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
			},
		},
		"split asymmetrical node with odd number of cells": {
			node: []byte{
				1, 7, 0, 24, 0, 30, 0, 38, 0, 44, 0, 50, 0, 56, 0, 64, 0,
				2, 0, 'a', 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b',
				2, 0, 'c', 'c', 2, 0, 'c', 'c', 1, 0, 'd', 1, 0, 'd',
				1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f',
				1, 0, 'g', 3, 0, 'g', 'g', 'g',
			},
			wantL: []byte{
				1, 3, 0, 16, 0, 22, 0, 30, 0,
				2, 0, 'a', 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b', 2, 0, 'c', 'c', 2, 0, 'c', 'c',
			},
			wantR: []byte{
				1, 4, 0, 17, 0, 23, 0, 29, 0, 37, 0,
				1, 0, 'd', 1, 0, 'd', 1, 0, 'e', 1, 0, 'e', 1, 0, 'f', 1, 0, 'f', 1, 0, 'g', 3, 0, 'g', 'g', 'g',
			},
		},
		"split asymmetrical node with large last cell": {
			node: []byte{
				1, 5, 0, 19, 0, 25, 0, 31, 0, 37, 0, 61, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b',
				1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
				10, 0, 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e',
				10, 0, 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e',
			},
			wantL: []byte{
				1, 4, 0, 17, 0, 23, 0, 29, 0, 35, 0,
				1, 0, 'a', 1, 0, 'a', 1, 0, 'b', 1, 0, 'b', 1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd',
			},
			wantR: []byte{
				1, 1, 0, 29, 0,
				10, 0, 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e',
				10, 0, 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e', 'e',
			},
		},
		"split asymmetrical node with large first cell": {
			node: []byte{
				1, 5, 0, 37, 0, 43, 0, 49, 0, 55, 0, 61, 0,
				10, 0, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a',
				10, 0, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a',
				1, 0, 'b', 1, 0, 'b', 1, 0, 'c', 1, 0, 'c',
				1, 0, 'd', 1, 0, 'd', 1, 0, 'e', 1, 0, 'e',
			},
			wantL: []byte{
				1, 1, 0, 29, 0,
				10, 0, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a',
				10, 0, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a',
			},
			wantR: []byte{
				1, 4, 0, 17, 0, 23, 0, 29, 0, 35, 0,
				1, 0, 'b', 1, 0, 'b', 1, 0, 'c', 1, 0, 'c', 1, 0, 'd', 1, 0, 'd', 1, 0, 'e', 1, 0, 'e',
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			gotL, gotR := Split(c.node)
			assert.Equal(t, c.wantL, gotL)
			assert.Equal(t, c.wantR, gotR)
		})
	}
}

// Fixtures

func fixtureNodeEvenNumberOfKeys() []byte {
	return []byte{
		1,
		10, 0,
		// Offsets
		33, 0,
		43, 0,
		53, 0,
		63, 0,
		73, 0,
		83, 0,
		93, 0,
		103, 0,
		113, 0,
		123, 0,
		// K-V Pairs
		3, 0, 'a', 'a', 'a', 3, 0, 'a', 'a', 'a',
		3, 0, 'b', 'b', 'b', 3, 0, 'b', 'b', 'b',
		3, 0, 'c', 'c', 'c', 3, 0, 'c', 'c', 'c',
		3, 0, 'd', 'd', 'd', 3, 0, 'd', 'd', 'd',
		3, 0, 'e', 'e', 'e', 3, 0, 'e', 'e', 'e',
		3, 0, 'f', 'f', 'f', 3, 0, 'f', 'f', 'f',
		3, 0, 'g', 'g', 'g', 3, 0, 'g', 'g', 'g',
		3, 0, 'h', 'h', 'h', 3, 0, 'h', 'h', 'h',
		3, 0, 'i', 'i', 'i', 3, 0, 'i', 'i', 'i',
		3, 0, 'j', 'j', 'j', 3, 0, 'j', 'j', 'j',
	}
}

func fixtureNodeOddNumberOfKeys() []byte {
	return []byte{
		1,
		11, 0,
		// Offsets
		35, 0,
		45, 0,
		55, 0,
		65, 0,
		75, 0,
		85, 0,
		95, 0,
		105, 0,
		115, 0,
		125, 0,
		135, 0,
		// K-V Pairs
		3, 0, 'a', 'a', 'a', 3, 0, 'a', 'a', 'a',
		3, 0, 'b', 'b', 'b', 3, 0, 'b', 'b', 'b',
		3, 0, 'c', 'c', 'c', 3, 0, 'c', 'c', 'c',
		3, 0, 'd', 'd', 'd', 3, 0, 'd', 'd', 'd',
		3, 0, 'e', 'e', 'e', 3, 0, 'e', 'e', 'e',
		3, 0, 'f', 'f', 'f', 3, 0, 'f', 'f', 'f',
		3, 0, 'g', 'g', 'g', 3, 0, 'g', 'g', 'g',
		3, 0, 'h', 'h', 'h', 3, 0, 'h', 'h', 'h',
		3, 0, 'i', 'i', 'i', 3, 0, 'i', 'i', 'i',
		3, 0, 'j', 'j', 'j', 3, 0, 'j', 'j', 'j',
		3, 0, 'k', 'k', 'k', 3, 0, 'k', 'k', 'k',
	}
}
